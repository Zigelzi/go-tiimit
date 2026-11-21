package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Zigelzi/go-tiimit/cmd/web/components"
	"github.com/Zigelzi/go-tiimit/internal/auth"
	"github.com/Zigelzi/go-tiimit/internal/db"
)

// The messages the sign-up form shows.
const (
	msgUsernameTaken = "That username is taken. Choose another one."
	msgSignUpFailed  = "Creating the account failed, try again later."
	// msgAlreadyLoggedIn is for a form sent from a tab opened before logging in.
	msgAlreadyLoggedIn = "You are already logged in. Log out to create a new account with this invitation."
	// msgLogInManually is for when the account exists but starting its session failed.
	msgLogInManually  = "Your account was created, but logging you in failed. Log in with your new username and password."
	msgUsernameEmpty  = "Choose a username."
	msgUsernameLength = "A username has 5 to 20 letters."
	msgUsernameChars  = "A username can have only letters (a to z, å, ä, ö)."
	msgPasswordShort  = "A password needs at least 10 characters."
	msgPasswordLong   = "A password can have at most 64 characters."
	msgPasswordBytes  = "This password is too long to store. Use a shorter one or fewer letters like ä and ö."
	msgPasswordSame   = "The password can't be the same as your username."
)

// inviteProblem is why an invite link cannot be used: a short title for the
// page that replaces the form, and a message that says what to do.
type inviteProblem struct {
	title, message string
}

// The ways an invitation can fail.
var (
	inviteUnknown = &inviteProblem{"Invitation doesn't work", "This invitation does not work. Ask the organizer for a new one."}
	inviteUsed    = &inviteProblem{"Invitation already used", "This invitation was already used. Ask the organizer for a new one."}
	inviteExpired = &inviteProblem{"Invitation expired", "This invitation has expired. Ask the organizer for a new one."}
)

// Errors createUserFromInvite returns when the account cannot be created for
// a reason the person can fix or should know about.
var (
	errInviteAlreadyUsed = errors.New("the invite was used or expired before the account was created")
	errUsernameTaken     = errors.New("the username is taken")
)

// handleSignUpPage shows the sign-up form for a working invite link, or why
// the link does not work. Someone who is already logged in gets a page with a
// log-out button instead of the form, after the link is checked. Opening the
// page never uses the link up, because chat apps open links to build a preview.
func (cfg *webConfig) handleSignUpPage(w http.ResponseWriter, r *http.Request) {
	_, problem, err := cfg.checkInvite(r.Context(), r.PathValue("token"), cfg.inviteNow())
	if err != nil {
		log.Printf("failed to check an invite: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		components.InviteLinkProblem("Something went wrong", msgSignUpFailed).Render(r.Context(), w)
		return
	}
	if problem != nil {
		components.InviteLinkProblem(problem.title, problem.message).Render(r.Context(), w)
		return
	}
	if user := auth.GetUserInfo(r.Context()); user.IsLoggedIn {
		components.AlreadyLoggedIn(user.Username).Render(r.Context(), w)
		return
	}
	components.SignUpPage(r.URL.Path).Render(r.Context(), w)
}

// handleSignUp creates an account from the invite link and logs the new user
// in. It checks the link first, then the username and password rules, then
// whether the username is taken, and hashes the password last, so someone
// without a working link learns nothing about usernames and cannot make the
// server hash passwords.
func (cfg *webConfig) handleSignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	now := cfg.inviteNow()

	inviteID, problem, err := cfg.checkInvite(ctx, r.PathValue("token"), now)
	if err != nil {
		signUpFailed(w, r, fmt.Errorf("checking the invite: %w", err))
		return
	}
	if problem != nil {
		renderFormError(w, r, problem.message)
		return
	}
	if auth.GetUserInfo(ctx).IsLoggedIn {
		renderFormError(w, r, msgAlreadyLoggedIn)
		return
	}

	username := auth.NormalizeUsername(r.FormValue("username"))
	password := r.FormValue("password")
	if problem := credentialsProblem(username, password); problem != "" {
		renderFormError(w, r, problem)
		return
	}

	_, err = cfg.queries.GetUserByUsername(ctx, username)
	if err == nil {
		renderFormError(w, r, msgUsernameTaken)
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		signUpFailed(w, r, fmt.Errorf("checking the username for invite %d: %w", inviteID, err))
		return
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		signUpFailed(w, r, fmt.Errorf("hashing the password for invite %d: %w", inviteID, err))
		return
	}

	userID, err := cfg.createUserFromInvite(ctx, inviteID, username, hashedPassword, now)
	switch {
	case errors.Is(err, errInviteAlreadyUsed):
		renderFormError(w, r, inviteUsed.message)
		return
	case errors.Is(err, errUsernameTaken):
		renderFormError(w, r, msgUsernameTaken)
		return
	case err != nil:
		signUpFailed(w, r, fmt.Errorf("creating the account from invite %d: %w", inviteID, err))
		return
	}
	log.Printf("invite %d created user %d", inviteID, userID)

	if err := cfg.startSession(w, r, userID); err != nil {
		// The account exists now, so "try again" would only meet a used link.
		log.Printf("failed to start a session for new user %d: %v", userID, err)
		renderFormError(w, r, msgLogInManually)
		return
	}
	w.Header().Set("HX-Redirect", "/")
}

// inviteNow returns the current time in UTC, cut to whole seconds, for
// storing and comparing invite times. SQLite compares the stored times as
// text, and a time with no fraction ("...26Z") would sort after one with a
// fraction in the same second ("...26.5Z"); whole seconds keep the order right.
func (cfg *webConfig) inviteNow() time.Time {
	return cfg.now().UTC().Truncate(time.Second)
}

// checkInvite looks up the invite with token at time now. It returns the
// invite's id and, when the link cannot be used, a message saying why. The
// error is set only when the lookup itself failed.
func (cfg *webConfig) checkInvite(ctx context.Context, token string, now time.Time) (int64, *inviteProblem, error) {
	status, err := cfg.queries.GetUserInviteStatus(ctx, db.GetUserInviteStatusParams{
		Now:   now,
		Token: token,
	})
	if errors.Is(err, sql.ErrNoRows) {
		return 0, inviteUnknown, nil
	}
	if err != nil {
		return 0, nil, err
	}
	if status.IsUsed == 1 {
		return status.ID, inviteUsed, nil
	}
	if status.IsExpired == 1 {
		return status.ID, inviteExpired, nil
	}
	return status.ID, nil, nil
}

// credentialsProblem returns the message for the first rule the normalized
// username or the password breaks, or "" when both follow the rules.
func credentialsProblem(username, password string) string {
	switch err := auth.ValidateUsername(username); {
	case errors.Is(err, auth.ErrUsernameEmpty):
		return msgUsernameEmpty
	case errors.Is(err, auth.ErrUsernameLength):
		return msgUsernameLength
	case errors.Is(err, auth.ErrUsernameLetters):
		return msgUsernameChars
	}
	switch err := auth.ValidatePassword(password, username); {
	case errors.Is(err, auth.ErrPasswordShort):
		return msgPasswordShort
	case errors.Is(err, auth.ErrPasswordLong):
		return msgPasswordLong
	case errors.Is(err, auth.ErrPasswordBytes):
		return msgPasswordBytes
	case errors.Is(err, auth.ErrPasswordSameAsUsername):
		return msgPasswordSame
	}
	return ""
}

// createUserFromInvite uses up the invite and creates the account in one
// transaction, so either both happen or neither does. It returns the new
// user's id. The invite is used up with one UPDATE that matches only an
// unused, unexpired invite, so of two people sending the form only one gets
// the account; the other gets errInviteAlreadyUsed. If both arrive within the
// same few milliseconds, the database has no busy timeout yet (tech review
// #7), so the second may get a plain error instead; still only one account is
// created. A taken username gives errUsernameTaken and leaves the invite
// working.
func (cfg *webConfig) createUserFromInvite(ctx context.Context, inviteID int64, username, hashedPassword string, now time.Time) (int64, error) {
	tx, err := cfg.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("starting the transaction: %w", err)
	}
	defer tx.Rollback()
	queries := cfg.queries.WithTx(tx)

	changed, err := queries.UseUserInvite(ctx, db.UseUserInviteParams{
		Now: sql.NullTime{Time: now, Valid: true},
		ID:  inviteID,
	})
	if err != nil {
		return 0, fmt.Errorf("using up the invite: %w", err)
	}
	if changed != 1 {
		return 0, errInviteAlreadyUsed
	}

	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Username:       username,
		HashedPassword: hashedPassword,
	})
	if db.IsUniqueViolation(err) {
		return 0, errUsernameTaken
	}
	if err != nil {
		return 0, fmt.Errorf("creating the user: %w", err)
	}

	err = queries.SetUserInviteUsedBy(ctx, db.SetUserInviteUsedByParams{
		UsedByUserID: sql.NullInt64{Int64: user.ID, Valid: true},
		ID:           inviteID,
	})
	if err != nil {
		return 0, fmt.Errorf("recording who used the invite: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("committing: %w", err)
	}
	return user.ID, nil
}

// renderFormError shows message in the sign-up form's error area.
func renderFormError(w http.ResponseWriter, r *http.Request, message string) {
	components.FormError(message).Render(r.Context(), w)
}

// signUpFailed logs err and shows the generic failure message in the form.
func signUpFailed(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("failed to sign up: %v", err)
	renderFormError(w, r, msgSignUpFailed)
}
