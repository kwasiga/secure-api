package model

import "time"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
	RoleGuest Role = "guest"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	LastName  string    `json:"last_name"`
	FirstName string    `json:"first_name"`
	Role      Role      `json:"role" db:"role"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
