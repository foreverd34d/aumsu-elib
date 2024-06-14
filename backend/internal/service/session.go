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

type sessionRepo interface {
	Create(ctx context.Context, userID int, input *model.NewToken) (*model.Token, error)
	GetUserFromSession(ctx context.Context, sessionID int) (*model.User, error)
	PopByRefreshToken(ctx context.Context, refreshToken string) (*model.Token, error)
	UpdateRefreshToken(ctx context.Context, sessionID int, update *model.NewToken) (*model.Token, error)
	EndSession(ctx context.Context, sessionID int) error
}

type SessionService struct {
	user userRepo
	session sessionRepo
	signingKey []byte
}

func NewSessionService(user userRepo, session sessionRepo, signingKey []byte) *SessionService {
	return &SessionService{
		user: user,
		session: session,
		signingKey: signingKey,
	}
}

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

func (ss *SessionService) Delete(ctx context.Context, refreshToken string) error {
	token, _ := ss.session.PopByRefreshToken(ctx, refreshToken)
	err := ss.session.EndSession(ctx, token.SessionID)
	if err != nil {
		return fmt.Errorf("end the session %v: %w", token.SessionID, err)
	}
	return err
}

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

func hashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

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

func createNewToken() *model.NewToken {
	return &model.NewToken{
		RefreshToken: generateRefreshToken(),
		ExpiresAt: int(time.Now().Add(30 * 24 * time.Hour).Unix()),
	}
}

func generateRefreshToken() string {
	buf := make([]byte, 32)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}
