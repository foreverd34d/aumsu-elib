package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Token представляет токен обновления,
// предназначеного для выдачи нового токена доступа.
type Token struct {
	ID           int    `json:"tokenID" db:"token_id"`           // номер
	RefreshToken string `json:"refreshToken" db:"refresh_token"` // токен обновления
	ExpiresAt    int    `json:"expiresAt" db:"expires_at"`       // время истечения срока действия токена, хранится в формате unix
	SessionID    int    `json:"sessionID" db:"session_id"`       // номер сессии
}

// Session представляет запись сессии пользователя.
type Session struct {
	ID          int        `json:"sessionID" db:"session_id"`      // номер
	LoggedInAt  time.Time  `json:"loggedInAt" db:"logged_in_at"`   // время входа в систему
	LoggedOutAt *time.Time `json:"loggedOutAt" db:"logged_out_at"` // время выхода из системы
	UserID      int        `json:"userID" db:"user_id"`            // номер пользователя
}

// JWTClaims представляет пользовательскую полезную нагрузку jwt токена.
type JWTClaims struct {
	jwt.RegisteredClaims
	Role UserRole
}

// UserRole представляет роль пользователя для использования в jwt токенах.
// Чем больше значение, тем выше полномочия
type UserRole int

const (
	StudentRole UserRole = iota // студент
	TeacherRole                 // преподаватель
	ManagerRole                 // руководитель
	AdminRole                   // администратор
)
