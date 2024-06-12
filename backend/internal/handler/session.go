package handler

import (
	"context"
	"libserver/internal/model"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type sessionService interface {
	Create(ctx context.Context, credentials *model.Credentials) (jwt string, session *model.Session, err error)
	Update(ctx context.Context, refreshToken string) (newjwt string, newSession *model.Session, err error)
	Delete(ctx context.Context, refreshToken string) error
}

type refreshTokenRequest struct {
	Token string `json:"refreshToken"`
}

func extractUserFromContext(c echo.Context) (*model.TokenClaims, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, echo.ErrUnauthorized
	}
	userClaims, ok := token.Claims.(*model.TokenClaims)
	if !ok {
		return nil, echo.ErrUnauthorized
	}
	return userClaims, nil
}

func (h *Handler) CreateSession(c echo.Context) error {
	credentials := new(model.Credentials)
	if err := c.Bind(credentials); err != nil {
		return echo.ErrBadRequest
	}
	jwt, session, err := h.Session.Create(c.Request().Context(), credentials)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"accessToken": jwt,
		"refreshToken": session.RefreshToken,
	})
}

func (h *Handler) UpdateSession(c echo.Context) error {
	var refreshToken refreshTokenRequest
	if err := c.Bind(&refreshToken); err != nil {
		return echo.ErrBadRequest
	}
	jwt, session, err := h.Session.Update(c.Request().Context(), refreshToken.Token)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"accessToken": jwt,
		"refreshToken": session.RefreshToken,
	})
}

func (h *Handler) DeleteSession(c echo.Context) error {
	var refreshToken refreshTokenRequest
	if err := c.Bind(&refreshToken); err != nil {
		return echo.ErrBadRequest
	}
	err := h.Session.Delete(c.Request().Context(), refreshToken.Token)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
