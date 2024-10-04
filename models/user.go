package models

import "time"

type UserRole string

const (
	USER  UserRole = "user"
	ADMIN UserRole = "admin"
)

func RoleStrConv(role string) UserRole {
	switch role {
	case "admin":
		return ADMIN
	default:
		return USER
	}
}

func (u UserRole) String() string {
	switch u {
	case ADMIN:
		return "admin"
	default:
		return "user"
	}
}

type UserResponse struct {
	ID        int      `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Role      UserRole `json:"role"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Role      UserRole  `json:"role"`
	Invoices  []Invoice `json:"invoices"`
	CreatedAt time.Time `json:"created_at"`
}
