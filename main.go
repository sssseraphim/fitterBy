package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sssseraphim/fitterBy/internal/auth"
	"github.com/sssseraphim/fitterBy/internal/database"
	"github.com/sssseraphim/fitterBy/internal/handlers"
	"github.com/sssseraphim/fitterBy/internal/middleware"
	"github.com/sssseraphim/fitterBy/internal/services"
)

type apiConfig struct {
	dbQueries          *database.Queries
	jwtTokenSecret     string
	refreshTokenSecret string
}

func main() {
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}
	var cfg apiConfig
	cfg.dbQueries = database.New(db)
	cfg.jwtTokenSecret = os.Getenv("JWT_SECRET")
	cfg.refreshTokenSecret = os.Getenv("REFRESH_SECRET")

	jwtConfig := &auth.JWTConfig{
		AccessTokenSecret:  cfg.jwtTokenSecret,
		RefreshTokenSecret: cfg.refreshTokenSecret,
		AccessTokenExpiry:  60 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}

	authHandler := &handlers.AuthHandler{
		DB:           cfg.dbQueries,
		TokenService: *services.NewTokenService(cfg.dbQueries, jwtConfig),
		JWTConfig:    jwtConfig,
	}
	mux := http.NewServeMux()
	// Serve static files (CSS, JS, images)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Auth endpoints
	mux.HandleFunc("POST /api/auth/signup", authHandler.HandleSignup)
	mux.HandleFunc("POST /api/auth/login", authHandler.HandleLogin)

	userHandler := &handlers.UserHandler{
		DB: cfg.dbQueries,
	}
	// User endpoints
	mux.HandleFunc("GET /api/users/{user_id}", userHandler.HandleGetUser)
	mux.Handle("GET /api/me", middleware.AuthMiddleware(jwtConfig)(http.HandlerFunc(userHandler.HandleGetCurrentUser)))
	mux.Handle("PATCH /api/me/bio", middleware.AuthMiddleware(jwtConfig)(http.HandlerFunc(userHandler.HandlerUpdateBio)))

	postHandler := &handlers.PostHandler{
		DB: cfg.dbQueries,
	}
	// Posts endpoint
	mux.HandleFunc("GET /api/posts/{post_id}", postHandler.HandleGetPost)
	mux.Handle("GET /api/posts/followed", middleware.AuthMiddleware(jwtConfig)(http.HandlerFunc(postHandler.HandleGetFollowedPosts)))
	mux.Handle("POST /api/posts", middleware.AuthMiddleware(jwtConfig)(http.HandlerFunc(postHandler.HandleCreatePost)))

	log.Println(" Servin from  http://localhost:8080/")

	server := &http.Server{Handler: mux, Addr: ":8080"}
	err = server.ListenAndServe()
	fmt.Println(err)
}

// serveTemplate serves HTML pages from templates folder
func serveTemplate(templateName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For now, just serve the HTML file
		// Later we'll add proper template rendering with data
		htmlPath := filepath.Join("static", "templates", templateName)

		htmlBytes, err := os.ReadFile(htmlPath)
		if err != nil {
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write(htmlBytes)
	}
}

// serveTrainerProfile handles dynamic trainer profile pages
func serveTrainerProfile(w http.ResponseWriter, r *http.Request) {
	// For now, serve a generic trainer profile template
	// Later we'll fetch actual trainer data
	htmlPath := filepath.Join("static", "templates", "trainer-single.html")

	htmlBytes, err := os.ReadFile(htmlPath)
	if err != nil {
		http.Error(w, "Trainer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(htmlBytes)
}

// API handlers (will implement functionality later)
func apiTrainersHandler(w http.ResponseWriter, r *http.Request) {
	// Return mock data for frontend development
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"trainers": [], "message": "API coming soon!"}`)
}

func apiSignupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "Signup API coming soon!"}`)
}

func apiLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "Login API coming soon!"}`)
}
