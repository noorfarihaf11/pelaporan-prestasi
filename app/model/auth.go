package model

import (
	"time"
	"github.com/google/uuid"
    "github.com/golang-jwt/jwt/v5" 
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

type LoginRequest struct { 
    Username string `json:"username"` 
    Password string `json:"password"` 
} 
type LoginResponse struct { 
	User  User   `json:"user"` 
	Token string `json:"token"` 
} 
type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}
type RegisterRequest struct {
    FullName string `json:"full_name" validate:"unique,required"`
    Username string `json:"username" validate:"unique,required,min=3,max=50"`
    Email    string `json:"email" validate:"unique,required,email"`
    Password string `json:"password" validate:"required,min=6"`
}
type RefreshClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}
