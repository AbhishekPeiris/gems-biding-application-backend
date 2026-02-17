package domain

import "time"

type UserRole string

const (
	RoleAdmin  UserRole = "ADMIN"
	RoleSeller UserRole = "SELLER"
	RoleBuyer  UserRole = "BUYER"
)

type User struct {
	ID        int64     `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // never expose
	Role      UserRole  `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
