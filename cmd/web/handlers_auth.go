package main

import (
	"database/sql"
	"errors"
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
	username := r.FormValue("username")
	dbUser, err := cfg.queries.GetUserByUsername(r.Context(), username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.Write([]byte("<p class='text-red-500'>Incorrect username or password </p>"))
			return
		} else {
			log.Printf("failed to get user [%s]: %v", username, err)
			w.Write([]byte("<p class='text-red-500'>Logging in failed, try again later.</p>"))
			return
		}
	}

	password := r.FormValue("password")
	err = auth.CheckPassword(password, dbUser.HashedPassword)
	if err != nil {
		w.Write([]byte("<p class='text-red-500'>Incorrect username or password </p>"))
		return
	}

	sessionToken, err := auth.MakeToken()
	if err != nil {
		log.Printf("failed to generate session token: %v", err)
		w.Write([]byte("<p class='text-red-500'>Logging in failed, try again later.</p>"))
		return
	}
	session, err := cfg.queries.StartUserSession(r.Context(), db.StartUserSessionParams{
		SessionID: sessionToken,
		UserID:    dbUser.ID,
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 14),
	})
	if err != nil {
		log.Printf("failed to start new session: %v", err)
		w.Write([]byte("<p class='text-red-500'>Logging in failed, try again later.</p>"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "user_session_id",
		Value:    session.SessionID,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
		Expires:  session.ExpiresAt,
	})
	w.Header().Add("HX-Redirect", "/")
}

func (cfg *webConfig) handleLogout(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("user_session_id")
	if err != nil {
		return
	}

	err = cfg.queries.EndUserSession(r.Context(), sessionCookie.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) == false {
			log.Printf("failed to log out session [%v]: %v", sessionCookie.Value, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	sessionCookie.Expires = time.Now().Add(-time.Hour)
	http.SetCookie(w, sessionCookie)
	w.Header().Add("HX-Redirect", "/")
}
