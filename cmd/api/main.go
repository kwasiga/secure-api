// main is the application entrypoint. It wires together config, database,
// middleware, and HTTP handlers, then starts the server.
package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/kwasiga/secure-api/db/sqlc"
	"github.com/kwasiga/secure-api/config"
	"github.com/kwasiga/secure-api/internal/handler"
	"github.com/kwasiga/secure-api/internal/middleware"
	"github.com/kwasiga/secure-api/internal/repository"
)

func main() {
	// Load configuration from environment (or .env file).
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// Connect to PostgreSQL using a connection pool.
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()

	// Build dependency tree: queries → repository → handlers.
	queries := db.New(pool)
	repo := repository.NewUserRepository(queries)

	authHandler := handler.NewAuthHandler(repo, cfg.JWTSecret)
	userHandler := handler.NewUserHandler(repo)
	adminHandler := handler.NewAdminHandler(repo)

	// Rate limiter: 10 req/s per IP, burst of 20.
	rl := middleware.NewRateLimiter(10, 20)

	r := chi.NewRouter()
	r.Use(rl.Limit) // applied globally to every route

	// Public routes — no authentication required.
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)

	// Protected routes — JWT required.
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.JWTSecret))

		r.Get("/me", userHandler.GetProfile)
		r.Put("/me", userHandler.UpdateProfile)

		// Admin-only routes — "admin" role required in addition to a valid JWT.
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireRole("admin"))
			r.Get("/admin/users", adminHandler.ListUsers)
		})
	})

	addr := ":" + cfg.Port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
