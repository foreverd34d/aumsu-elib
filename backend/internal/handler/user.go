package handler

import (
	"context"
	"libserver/internal/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)


type userService interface {
	Create(ctx context.Context, input *model.NewUser) (*model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
	Get(ctx context.Context, ID int) (*model.User, error)
	Update(ctx context.Context, ID int, update *model.NewUser) (*model.User, error)
	Delete(ctx context.Context, ID int) error
}

func (h *Handler) CreateUser(c echo.Context) error {
	userClaims, err := extractUserFromContext(c)
	if err != nil {
		return echo.ErrUnauthorized
	}
	if userClaims.Role < model.AdminRole {
		return echo.ErrForbidden
	}

	newUser := new(model.NewUser)
	if err := c.Bind(newUser); err != nil {
		return echo.ErrBadRequest
	}
	user, err := h.User.Create(c.Request().Context(), newUser)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"id": user.ID,
	})
}

func (h *Handler) GetAllUsers(c echo.Context) error {
	userClaims, err := extractUserFromContext(c)
	if err != nil {
		return echo.ErrUnauthorized
	}
	if userClaims.Role < model.AdminRole {
		return echo.ErrForbidden
	}

	users, err := h.User.GetAll(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUser(c echo.Context) error {
	userClaims, err := extractUserFromContext(c)
	if err != nil {
		return echo.ErrUnauthorized
	}
	if userClaims.Role < model.AdminRole {
		return echo.ErrForbidden
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	}
	user, err := h.User.Get(c.Request().Context(), userID)
	if err != nil {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(c echo.Context) error {
	userClaims, err := extractUserFromContext(c)
	if err != nil {
		return echo.ErrUnauthorized
	}
	if userClaims.Role < model.AdminRole {
		return echo.ErrForbidden
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	}
	userUpdate := new(model.NewUser)
	if err := c.Bind(userUpdate); err != nil {
		return echo.ErrBadRequest
	}

	_, err = h.User.Update(c.Request().Context(), userID, userUpdate)
	if err != nil {
		return echo.ErrNotFound
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DeleteUser(c echo.Context) error {
	userClaims, err := extractUserFromContext(c)
	if err != nil {
		return echo.ErrUnauthorized
	}
	if userClaims.Role < model.AdminRole {
		return echo.ErrForbidden
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	}
	err = h.User.Delete(c.Request().Context(), userID)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
