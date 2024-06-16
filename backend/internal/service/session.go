package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/foreverd34d/aumsu-elib/internal/errs"
	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

// SessionRepo определяет методы хранилища токенов и сессий.
type SessionRepo interface {
	// Create создает новый токен обновления, записывает время начала сессии для пользователя
	// и возвращает токен обновления с номером или ошибку.
	Create(ctx context.Context, userID int, input *model.NewToken) (*model.Token, error)

	// GetUserFromSession возвращает пользователя по номеру его сессии или ошибку.
	GetUserFromSession(ctx context.Context, sessionID int) (*model.User, error)

	// PopByRefreshToken удаляет токен обновления из базы данных и возвращает всю информацию о нем или ошибку.
	// Если токен не был найден, то возвращается ошибка [errs.NotFound].
	PopByRefreshToken(ctx context.Context, refreshToken string) (*model.Token, error)

	// UpdateRefreshToken создает новый токен обновления для сессии и возвращает токен с номером или ошибку.
	UpdateRefreshToken(ctx context.Context, sessionID int, update *model.NewToken) (*model.Token, error)

	// EndSession записывает время окончания сессии и возвращает ошибку, если таковая есть.
	EndSession(ctx context.Context, sessionID int) error
}

// SessionService реализует методы для работы с токенами и сессиями
// и реализует интерфейс [handler.SessionService].
type SessionService struct {
	user UserRepo
	session SessionRepo
	signingKey []byte
}

// NewSessionService возвращает новый экземпляр [SessionService].
func NewSessionService(user UserRepo, session SessionRepo, signingKey []byte) *SessionService {
	return &SessionService{
		user: user,
		session: session,
		signingKey: signingKey,
	}
}

// Create создает пару из jwt токена и токена обновления и записывает время начала сессия пользователя.
// Если имя пользователя не найдено или пароль не совпадает с сохраненным,
// то возвращается ошибка [errs.InvalidLogin] или [errs.InvalidPassword].
func (ss *SessionService) Create(ctx context.Context, credentials *model.Credentials) (jwt string, refreshToken *model.Token, err error) {
	dbCredentials, err := ss.user.GetCredentialsByLogin(ctx, credentials.Username)
	if err != nil {
		err = fmt.Errorf("get the user %s: %w", credentials.Username, err)
		return
	}

	if dbCredentials.PasswordHash != hashPassword(credentials.Password) {
		err = fmt.Errorf("create a session: %w", errs.InvalidPassword)
		return
	}

	role, err := ss.user.GetRole(ctx, dbCredentials.UserID)
	if err != nil {
		err = fmt.Errorf("get user role: %w", err)
		return
	}

	refreshToken, err = ss.session.Create(ctx, dbCredentials.UserID, createNewToken())
	if err != nil {
		err = fmt.Errorf("create new session: %w", err)
		return
	}

	jwt, err = createJWT(dbCredentials.UserID, role, ss.signingKey)
	if err != nil {
		err = fmt.Errorf("create a jwt token: %w", err)
	}
	return
}

// Update создает новую пару токенов по токену обновления. Сессия при этом не кончается,
// а старый токен обновления становится невалидным.
// Если токен обновления истек, то возвращается ошибка [errs.RefreshExpired].
func (ss *SessionService) Update(ctx context.Context, refreshToken string) (newjwt string, newRefreshToken *model.Token, err error) {
	token, err := ss.session.PopByRefreshToken(ctx, refreshToken)
	if err != nil {
		err = fmt.Errorf("pop the refresh token %v: %w", refreshToken, err)
		return
	}

	if token.ExpiresAt < int(time.Now().Unix()) {
		ss.session.EndSession(ctx, token.SessionID)
		err = fmt.Errorf("update the session %v: %w", token.SessionID, errs.RefreshExpired)
		return
	}

	newRefreshToken, err = ss.session.UpdateRefreshToken(ctx, token.SessionID, createNewToken())
	if err != nil {
		ss.session.EndSession(ctx, token.SessionID)
		err = fmt.Errorf("update the refresh token %v: %w", token.ID, err)
		return
	}

	user, err := ss.session.GetUserFromSession(ctx, token.SessionID)
	if err != nil {
		ss.session.EndSession(ctx, token.SessionID)
		err = fmt.Errorf("get user from session %v: %w", token.SessionID, err)
		return
	}

	role, err := ss.user.GetRole(ctx, user.ID)
	if err != nil {
		ss.session.EndSession(ctx, token.SessionID)
		err = fmt.Errorf("get a role for user %v: %w", user.ID, err)
		return
	}

	newjwt, err = createJWT(user.ID, role, ss.signingKey)
	if err != nil {
		ss.session.EndSession(ctx, token.SessionID)
		err = fmt.Errorf("create the JWT with userID %v: %w", user.ID, err)
	}
	return
}

// Delete делает токен обновления невалидным и записывает время окончания сессии.
func (ss *SessionService) Delete(ctx context.Context, refreshToken string) error {
	token, _ := ss.session.PopByRefreshToken(ctx, refreshToken)
	err := ss.session.EndSession(ctx, token.SessionID)
	if err != nil {
		return fmt.Errorf("end the session %v: %w", token.SessionID, err)
	}
	return err
}

// getRoleFromName возвращает роль исходя из ее названия.
func getRoleFromName(roleName string) model.UserRole {
	var role model.UserRole
	switch roleName {
	case "student":
		role = model.StudentRole
	case "teacher":
		role = model.TeacherRole
	case "manager":
		role = model.ManagerRole
	case "admin":
		role = model.AdminRole
	}
	return role
}

// hashPassword хэширует пароль алгоритмом sha-256.
func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

// createJWT создает новый jwt токен пользователя с его ролью
// и подписывает его при помощи переданого ключа алгоритмом sha-256.
func createJWT(userID int, roleName string, signingKey []byte) (string, error) {
	role := getRoleFromName(roleName)
	claims := &model.JWTClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: strconv.Itoa(userID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signingKey)
}

// createNewToken создает новый токен обновления со сроком действия в месяц.
func createNewToken() *model.NewToken {
	return &model.NewToken{
		RefreshToken: generateRefreshToken(),
		ExpiresAt: int(time.Now().Add(30 * 24 * time.Hour).Unix()),
	}
}

// generateRefreshToken возвращает случайно сгенерированную строку
// длиной 32 символа, представляющую из себя токен обновления.
func generateRefreshToken() string {
	buf := make([]byte, 32)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}
