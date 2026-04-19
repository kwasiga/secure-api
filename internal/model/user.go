package model

import (
	"time"

	db "github.com/kwasiga/secure-api/db/sqlc"
)

type User struct {
	ID        int32       `json:"id"`
	Email     string      `json:"email"`
	LastName  string      `json:"last_name"`
	FirstName string      `json:"first_name"`
	Role      db.UserRole `json:"role" db:"role"`
	Password  string      `json:"-"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
}
