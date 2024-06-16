package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/labstack/echo/v4"
)

// SessionService определяет методы для работы с токенами и сессиями.
type SessionService interface {
	// Create создает пару из jwt токена и токена обновления и записывает время начала сессия пользователя.
	// Если имя пользователя не найдено или пароль не совпадает с сохраненным,
	// то возвращается ошибка [errs.InvalidLogin] или [errs.InvalidPassword].
	Create(ctx context.Context, credentials *model.Credentials) (jwt string, token *model.Token, err error)

	// Update создает новую пару токенов по токену обновления. Сессия при этом не кончается,
	// а старый токен обновления становится невалидным.
	// Если токен обновления истек, то возвращается ошибка [errs.RefreshExpired].
	Update(ctx context.Context, refreshToken string) (newjwt string, newToken *model.Token, err error)

	// Delete делает токен обновления невалидным и записывает время окончания сессии.
	Delete(ctx context.Context, refreshToken string) error
}

// refreshTokenRequest оборачивает токен обновления в json-объект для получения из тела запроса.
type refreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validator:"required"` // токен обновления
}

// CreateSession получает данные для входа из тела запроса и создает пару из jwt токена
// и токена обновления. Если имя пользователя не найдено или пароль не совпадает с сохраненным,
// то возвращается ошибка [erss.InvalidLogin] или [errs.InvalidPassword]. В ответе возвращаются
// jwt токен и токен обновления.
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

// UpdateSession получает токен обновления из тела запроса и создает новую пару токенов.
// Сессия при этом не кончается, а старый токен обновления становится невалидным.
// Если токен обновления истек, то возвращается ошибка [errs.RefreshExpired].
// В ответе возвращаются jwt токен и токен обновления.
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

// DeleteSession делает токен обновления невалидным и записывает время окончания сессии.
// В ответе ничего не возвращается.
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
