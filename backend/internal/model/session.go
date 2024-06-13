package model

import "github.com/golang-jwt/jwt/v5"

type Session struct {
	ID           int    `json:"sessionID" db:"session_id"`
	RefreshToken string `json:"refreshToken" db:"refresh_token"`
	ExpiresAt    int    `json:"expiresAt" db:"expires_at"`
	UserID       int    `json:"userID" db:"user_id"`
}

type TokenClaims struct {
	jwt.RegisteredClaims
	Role UserRole
}

type UserRole int

const (
	StudentRole UserRole = iota
	TeacherRole
	ManagerRole
	AdminRole
)
