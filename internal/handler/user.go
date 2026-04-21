package handler

import (
	"encoding/json"
	"net/http"

	db "github.com/kwasiga/secure-api/db/sqlc"
	"github.com/kwasiga/secure-api/internal/middleware"
	"github.com/kwasiga/secure-api/internal/repository"
	"github.com/kwasiga/secure-api/internal/validator"
)

// UserHandler handles authenticated user profile routes.
type UserHandler struct {
	repo *repository.UserRepository
}

// NewUserHandler constructs a UserHandler with the given repository.
func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// GetProfile handles GET /me.
// Returns the authenticated user's profile as JSON.
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.repo.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateProfileRequest is the expected body for PUT /me.
type UpdateProfileRequest struct {
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
}

// UpdateProfile handles PUT /me.
// Validates the request body, updates the authenticated user's name, and returns the updated record.
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req UpdateProfileRequest
	if err := validator.Validate(w, r, &req); err != nil {
		return
	}

	user, err := h.repo.UpdateUser(r.Context(), db.UpdateUserParams{
		ID:        claims.UserID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	})
	if err != nil {
		http.Error(w, "could not update user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
