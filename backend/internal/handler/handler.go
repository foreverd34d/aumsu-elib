package handler

import (
	"libserver/internal/model"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	User userService
	Session sessionService
	Group groupService
}

func checkRole(c echo.Context, role model.UserRole) error {
	user, err := extractUserFromContext(c)
	if err != nil {
		return echo.ErrUnauthorized
	}
	if user.Role < role {
		return echo.ErrForbidden
	}
	return nil
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

