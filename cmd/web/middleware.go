package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/Zigelzi/go-tiimit/internal/auth"
)

func (cfg *webConfig) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Skip authenticating static assets (CSS, JS, etc.)
		if strings.HasPrefix(r.URL.Path, "/static/") {
			next.ServeHTTP(w, r)
			return
		}
		userInfo := cfg.getUserInfoFromRequest(r)
		ctxWithUserInfo := auth.WithUserInfo(r.Context(), userInfo)
		next.ServeHTTP(w, r.WithContext(ctxWithUserInfo))
	})
}

func (cfg *webConfig) getUserInfoFromRequest(r *http.Request) auth.UserInfo {
	sessionCookie, err := r.Cookie("user_session_id")
	if err != nil {
		return auth.UserInfo{
			IsLoggedIn: false,
		}
	}
	_, err = cfg.queries.GetActiveSession(r.Context(), sessionCookie.Value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) == false {
			log.Printf("failed to get active user session: %v", err)
		}
		return auth.UserInfo{IsLoggedIn: false}
	}

	return auth.UserInfo{
		IsLoggedIn: true,
	}
}

func requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo := auth.GetUserInfo(r.Context())

		if userInfo.IsLoggedIn == false {
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/login")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
