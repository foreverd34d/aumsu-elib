package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	ID           int    `json:"tokenID" db:"token_id"`
	RefreshToken string `json:"refreshToken" db:"refresh_token"`
	ExpiresAt    int    `json:"expiresAt" db:"expires_at"`
	SessionID    int    `json:"sessionID" db:"session_id"`
}

type Session struct {
	ID          int        `json:"sessionID" db:"session_id"`
	LoggedInAt  time.Time  `json:"loggedInAt" db:"logged_in_at"`
	LoggedOutAt *time.Time `json:"loggedOutAt" db:"logged_out_at"`
	UserID      int        `json:"userID" db:"user_id"`
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
