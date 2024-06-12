package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"libserver/internal/model"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)

type sessionRepo interface {
	Create(ctx context.Context, input *model.NewSession) (*model.Session, error)
	PopByRefreshToken(ctx context.Context, refreshToken string) (*model.Session, error)
}

type SessionService struct {
	userRepo userRepo
	sessionRepo sessionRepo
	signingKey []byte
}

func NewSessionService(user userRepo, session sessionRepo, signingKey []byte) *SessionService {
	return &SessionService{
		userRepo: user,
		sessionRepo: session,
		signingKey: signingKey,
	}
}

func (ss *SessionService) Create(ctx context.Context, credentials *model.Credentials) (jwt string, session *model.Session, err error) {
	user, err := ss.userRepo.GetByLogin(ctx, credentials.Username)
	if err != nil {
		return
	}

	if user.PasswordHash != hashPassword(credentials.Password) {
		err = ErrInvalidPassword
		return
	}

	role, err := ss.userRepo.GetRole(ctx, user.ID)
	if err != nil {
		return
	}

	session, err = ss.sessionRepo.Create(ctx, createNewSession(user.ID))
	if err != nil {
		return
	}

	jwt, err = createJWT(user.ID, role, ss.signingKey)
	return
}

func (ss *SessionService) Update(ctx context.Context, refreshToken string) (newjwt string, newSession *model.Session, err error) {
	session, err := ss.sessionRepo.PopByRefreshToken(ctx, refreshToken)
	if err != nil {
		return
	}

	if session.ExpiresAt > int(time.Now().Unix()) {
		err = ErrRefreshTokenExpired
		return
	}

	newSession, err = ss.sessionRepo.Create(ctx, createNewSession(session.UserID))
	if err != nil {
		return
	}

	role, err := ss.userRepo.GetRole(ctx, session.UserID)
	if err != nil {
		return
	}

	newjwt, err = createJWT(session.UserID, role, ss.signingKey)
	return
}

func (ss *SessionService) Delete(ctx context.Context, refreshToken string) error {
	_, err := ss.sessionRepo.PopByRefreshToken(ctx, refreshToken)
	return err
}

func getRoleFromName(roleName string) model.UserRole {
	var role model.UserRole
	switch roleName {
	case "student":
		role = model.StudentRole
	case "teacher":
		role = model.TeacherRole
	case "admin":
		role = model.AdminRole
	}
	return role
}

func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

func createJWT(userID int, roleName string, signingKey []byte) (string, error) {
	role := getRoleFromName(roleName)
	claims := &model.TokenClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: strconv.Itoa(userID),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signingKey)
}

func createNewSession(userID int) *model.NewSession {
	return &model.NewSession{
		RefreshToken: generateRefreshToken(),
		UserID: userID,
		ExpiresAt: int(time.Now().Add(30 * 24 * time.Hour).Unix()),
	}
}

func generateRefreshToken() string {
	buf := make([]byte, 32)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}
