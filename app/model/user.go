package model

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	RoleID       uuid.UUID `json:"role_id"`
	IsActive     bool      `json:"is_active"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type CreateUser struct {
    FullName string `json:"full_name" validate:"unique,required"`
    Username string `json:"username" validate:"unique,required,min=3,max=50"`
    Email    string `json:"email" validate:"unique,required,email"`
    Password string `json:"password" validate:"required,min=6"`
	RoleID   string `json:"role_id"`
}