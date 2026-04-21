package handler

import (
	"encoding/json"
	"net/http"

	"github.com/kwasiga/secure-api/internal/repository"
)

// AdminHandler handles admin-only routes.
type AdminHandler struct {
	repo *repository.UserRepository
}

// NewAdminHandler constructs an AdminHandler with the given repository.
func NewAdminHandler(repo *repository.UserRepository) *AdminHandler {
	return &AdminHandler{repo: repo}
}

// ListUsers handles GET /admin/users.
// Returns all registered users as a JSON array. Requires the "admin" role.
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.ListUsers(r.Context())
	if err != nil {
		http.Error(w, "could not fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
