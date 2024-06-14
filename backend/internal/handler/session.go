package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/labstack/echo/v4"
)

type sessionService interface {
	Create(ctx context.Context, credentials *model.Credentials) (jwt string, token *model.Token, err error)
	Update(ctx context.Context, refreshToken string) (newjwt string, newToken *model.Token, err error)
	Delete(ctx context.Context, refreshToken string) error
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (h *Handler) CreateSession(c echo.Context) error {
	credentials := new(model.Credentials)
	if err := bindAndValidate(c, credentials); err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("bind credentials: %w", err))
	}
	jwt, token, err := h.Session.Create(c.Request().Context(), credentials)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"accessToken": jwt,
		"refreshToken": token.RefreshToken,
	})
}

func (h *Handler) UpdateSession(c echo.Context) error {
	var refreshToken refreshTokenRequest
	if err := bindAndValidate(c, &refreshToken); err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("bind refresh token: %w", err))
	}
	jwt, token, err := h.Session.Update(c.Request().Context(), refreshToken.RefreshToken)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{
		"accessToken": jwt,
		"refreshToken": token.RefreshToken,
	})
}

func (h *Handler) DeleteSession(c echo.Context) error {
	var refreshToken refreshTokenRequest
	if err := bindAndValidate(c, &refreshToken); err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("bind refresh token: %w", err))
	}
	err := h.Session.Delete(c.Request().Context(), refreshToken.RefreshToken)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
