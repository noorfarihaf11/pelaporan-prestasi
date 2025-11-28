package model

import (
	_"time"
	"github.com/google/uuid"
    "github.com/golang-jwt/jwt/v5" 
)

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
	RoleID   uuid.UUID `json:"role_id"`
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
