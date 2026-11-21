package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Zigelzi/go-tiimit/cmd/web/components"
	"github.com/Zigelzi/go-tiimit/internal/auth"
	"github.com/Zigelzi/go-tiimit/internal/db"
)

// inviteLifetime is how long an invite link works.
const inviteLifetime = 48 * time.Hour

func (cfg *webConfig) handleViewUsersPage(w http.ResponseWriter, r *http.Request) {
	usernames, err := cfg.queries.ListUsernames(r.Context())
	if err != nil {
		log.Printf("failed to list users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	components.UsersPage(usernames).Render(r.Context(), w)
}

// handleCreateInvite stores a new invite and replies with its path and expiry
// time. The browser builds the full link from its own address.
func (cfg *webConfig) handleCreateInvite(w http.ResponseWriter, r *http.Request) {
	token, err := auth.MakeToken()
	if err != nil {
		log.Printf("failed to make an invite token: %v", err)
		renderError(w, r, "Couldn't create the invitation. Try again later.")
		return
	}
	now := cfg.inviteNow()
	expiresAt := now.Add(inviteLifetime)
	invite, err := cfg.queries.CreateUserInvite(r.Context(), db.CreateUserInviteParams{
		Token:           token,
		CreatedByUserID: auth.GetUserInfo(r.Context()).ID,
		CreatedAt:       now,
		ExpiresAt:       expiresAt,
	})
	if err != nil {
		log.Printf("failed to create an invite: %v", err)
		renderError(w, r, "Couldn't create the invitation. Try again later.")
		return
	}
	log.Printf("created invite %d", invite.ID)

	expiryText, err := formatFinnishTime(expiresAt)
	if err != nil {
		log.Printf("failed to format the expiry time of invite %d: %v", invite.ID, err)
		renderError(w, r, "Couldn't create the invitation. Try again later.")
		return
	}
	renderOK(w, r, components.InviteResult("/invites/"+token, expiryText))
}

// formatFinnishTime returns t in Finnish time, such as "Mon 23 Nov 18:40".
func formatFinnishTime(t time.Time) (string, error) {
	helsinki, err := time.LoadLocation("Europe/Helsinki")
	if err != nil {
		return "", err
	}
	return t.In(helsinki).Format("Mon 2 Jan 15:04"), nil
}
