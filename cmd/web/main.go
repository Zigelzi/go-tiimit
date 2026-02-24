package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/Zigelzi/go-tiimit/internal/db"
	"github.com/joho/godotenv"
)

//go:embed static
var staticFiles embed.FS

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load .env file: %v", err)
	}
	newDb, err := db.InitDB()
	if err != nil {
		log.Fatalf("initializing database failed: %v", err)
		return
	}
	defer newDb.Close()

	err = db.RunMigrations(newDb)
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
		return
	}

	cfg := webConfig{
		queries: db.New(newDb),
		db:      newDb,
		address: ":8080",
		env:     "development",
	}

	mux := http.NewServeMux()

	staticFs, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("failed to initialize the static files")
		return
	}

	fileserver := http.FileServer(http.FS(staticFs))
	mux.Handle("/static/", disableCacheInDevMode(http.StripPrefix("/static/", fileserver), cfg.env))

	mux.HandleFunc("/", cfg.handleIndexPage)

	// Practices
	mux.Handle("GET /practice/set-up", requireAuth(http.HandlerFunc(cfg.handleSetupPracticePage)))
	mux.Handle("GET /practice/{id}", requireAuth(http.HandlerFunc(cfg.handleViewPractice)))
	mux.Handle("POST /practice", requireAuth(http.HandlerFunc(cfg.handleCreatePractice)))

	// Auth
	mux.HandleFunc("GET /login", cfg.handleLoginPage)
	mux.HandleFunc("POST /login", cfg.handleLoginUser)
	mux.HandleFunc("POST /logout", cfg.handleLogout)

	muxWithAuth := cfg.authMiddleware(mux)
	server := http.Server{
		Handler: muxWithAuth,
		Addr:    cfg.address,
	}
	log.Printf("Starting server on address %s", cfg.address)
	server.ListenAndServe()
}

func disableCacheInDevMode(next http.Handler, env string) http.Handler {
	if env != "development" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
