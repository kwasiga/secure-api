// Package repository provides database access methods that map sqlc-generated
// types to the app-layer model types.
package repository

import (
	"context"

	db "github.com/kwasiga/secure-api/db/sqlc"
	"github.com/kwasiga/secure-api/internal/model"
)

// toModel converts a sqlc-generated User row into the app-layer model.User.
func toModel(u db.User) model.User {
	return model.User{
		ID:        u.ID,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		Password:  u.Password,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
	}
}

// UserRepository wraps the sqlc Queries and exposes domain-level operations.
type UserRepository struct {
	queries *db.Queries
}

// NewUserRepository constructs a UserRepository backed by the given Queries.
func NewUserRepository(queries *db.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

// CreateUser inserts a new user row and returns the created record.
func (r *UserRepository) CreateUser(ctx context.Context, arg db.CreateUserParams) (model.User, error) {
	u, err := r.queries.CreateUser(ctx, arg)
	if err != nil {
		return model.User{}, err
	}
	return toModel(u), nil
}

// GetUserByEmail fetches a user by email address.
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	u, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, err
	}
	return toModel(u), nil
}

// GetUserByID fetches a user by primary key.
func (r *UserRepository) GetUserByID(ctx context.Context, ID int32) (model.User, error) {
	u, err := r.queries.GetUserByID(ctx, ID)
	if err != nil {
		return model.User{}, err
	}
	return toModel(u), nil
}

// ListUsers returns all users ordered by ID.
func (r *UserRepository) ListUsers(ctx context.Context) ([]model.User, error) {
	users, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]model.User, len(users))
	for i, u := range users {
		result[i] = toModel(u)
	}
	return result, nil
}

// UpdateUser updates a user's first and last name and returns the updated record.
func (r *UserRepository) UpdateUser(ctx context.Context, arg db.UpdateUserParams) (model.User, error) {
	u, err := r.queries.UpdateUser(ctx, arg)
	if err != nil {
		return model.User{}, err
	}
	return toModel(u), nil
}
