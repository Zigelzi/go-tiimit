package main

import (
	"net/http"
	"testing"
)

// The router serves the embedded static files under /static/, which every
// page's styles and scripts depend on.
func TestStaticFilesAreServed(t *testing.T) {
	client := newTestClient(t, newTestApp(t))

	rec := client.get("/static/tabs.js")

	wantStatus(t, rec, http.StatusOK)
}
