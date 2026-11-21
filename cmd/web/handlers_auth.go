package main

import (
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

func (cfg *webConfig) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	component := components.LoginPage()
	component.Render(r.Context(), w)
}

func (cfg *webConfig) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	username := auth.NormalizeUsername(r.FormValue("username"))
	dbUser, err := cfg.queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			component := components.FormError("Incorrect username or password")
			component.Render(r.Context(), w)
			return
		} else {
			log.Printf("failed to get user [%s]: %v", username, err)
			component := components.FormError("Incorrect username or password")
			component.Render(r.Context(), w)
			return
		}
	}

	password := r.FormValue("password")
	err = auth.CheckPassword(password, dbUser.HashedPassword)
	if err != nil {
		component := components.FormError("Incorrect username or password")
		component.Render(r.Context(), w)
		return
	}

	if err := cfg.startSession(w, r, dbUser.ID); err != nil {
		log.Printf("failed to start a session for user %d: %v", dbUser.ID, err)
		component := components.FormError("Logging in failed, try again later.")
		component.Render(r.Context(), w)
		return
	}
	w.Header().Add("HX-Redirect", "/")
}

// sessionLifetime is how long a login lasts.
const sessionLifetime = 14 * 24 * time.Hour

// startSession logs userID in: it stores a new session and sets its cookie on
// w. Logging in and signing up both use it, so a session is made one way.
func (cfg *webConfig) startSession(w http.ResponseWriter, r *http.Request, userID int64) error {
	sessionToken, err := auth.MakeToken()
	if err != nil {
		return fmt.Errorf("generating the session token: %w", err)
	}
	session, err := cfg.queries.StartUserSession(r.Context(), db.StartUserSessionParams{
		SessionID: sessionToken,
		UserID:    userID,
		ExpiresAt: time.Now().UTC().Add(sessionLifetime),
	})
	if err != nil {
		return fmt.Errorf("storing the session: %w", err)
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "user_session_id",
		Value:    session.SessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
		Expires:  session.ExpiresAt,
	})
	return nil
}

func (cfg *webConfig) handleLogout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("user_session_id")
	if err != nil {
		// Nothing to end, but still leave the page, so the button never
		// seems to do nothing.
		w.Header().Add("HX-Redirect", "/")
		return
	}

	err = cfg.queries.EndUserSession(r.Context(), sessionCookie.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) == false {
			log.Printf("failed to log out session: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// A cookie read from the request has no path, so name it again; without
	// it the browser would not match the cookie it is meant to remove.
	sessionCookie.Path = "/"
	sessionCookie.Expires = time.Now().Add(-time.Hour)
	http.SetCookie(w, sessionCookie)
	w.Header().Add("HX-Redirect", "/")
}
