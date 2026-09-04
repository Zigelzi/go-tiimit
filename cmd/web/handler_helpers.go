package main

import (
	"net/http"

	"github.com/Zigelzi/go-tiimit/cmd/web/components"
	"github.com/a-h/templ"
)

func renderError(w http.ResponseWriter, r *http.Request, msg string) {
	w.Header().Set("HX-Reswap", "none")
	components.ErrorBanner(msg).Render(r.Context(), w)
}

func renderOK(w http.ResponseWriter, r *http.Request, component templ.Component) {
	component.Render(r.Context(), w)
	components.ErrorBanner("").Render(r.Context(), w)
}
