// Package handler contains HTTP handlers for all API routes.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kwasiga/secure-api/internal/auth"
	db "github.com/kwasiga/secure-api/db/sqlc"
	"github.com/kwasiga/secure-api/internal/repository"
	"github.com/kwasiga/secure-api/internal/validator"
)

// AuthHandler handles public authentication routes.
type AuthHandler struct {
	jwtSecret string
	repo      *repository.UserRepository
}

// NewAuthHandler constructs an AuthHandler with the given repository and JWT secret.
func NewAuthHandler(repo *repository.UserRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{repo: repo, jwtSecret: jwtSecret}
}

// RegisterRequest is the expected body for POST /register.
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	Password  string `json:"password" validate:"required,min=8"`
}

// LoginRequest is the expected body for POST /login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Register handles POST /register.
// Validates the request, hashes the password, creates the user, and returns the new record.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := validator.Validate(w, r, &req); err != nil {
		return
	}

	hashed, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	user, err := h.repo.CreateUser(r.Context(), db.CreateUserParams{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  hashed,
		Role:      db.UserRoleUser,
	})
	if err != nil {
		http.Error(w, "could not create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Login handles POST /login.
// Validates credentials and returns a short-lived access token and a long-lived refresh token.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := validator.Validate(w, r, &req); err != nil {
		return
	}

	user, err := h.repo.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := auth.CheckPassword(req.Password, user.Password); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := auth.GenerateTokens(user.ID, string(user.Role), h.jwtSecret)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
