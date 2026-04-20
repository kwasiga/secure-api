package repository

import (
	"context"

	db "github.com/kwasiga/secure-api/db/sqlc"
	"github.com/kwasiga/secure-api/internal/model"
)

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

type UserRepository struct {
	queries *db.Queries
}

func NewUserRepository(queries *db.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) CreateUser(ctx context.Context, arg db.CreateUserParams) (model.User, error) {
	u, err := r.queries.CreateUser(ctx, arg)
	if err != nil {
		return model.User{}, err
	}
	return toModel(u), nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	u, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, err
	}
	return toModel(u), nil
}
func (r *UserRepository) GetUserByID(ctx context.Context, ID int32) (model.User, error) {
	u, err := r.queries.GetUserByID(ctx, ID)
	if err != nil {
		return model.User{}, err
	}
	return toModel(u), nil
}
func (r *UserRepository) UpdateUser(ctx context.Context, arg db.UpdateUserParams) (model.User, error) {
	u, err := r.queries.UpdateUser(ctx, arg)
	if err != nil {
		return model.User{}, err
	}
	return toModel(u), nil
}
