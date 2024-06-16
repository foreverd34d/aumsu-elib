package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/foreverd34d/aumsu-elib/internal/errs"
	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/labstack/echo/v4"
)

// UserService декларирует методы для работы со пользователями и данными для их входа.
type UserService interface {
	// Create создает нового пользователя и его данные для входа и возвращает пользователя с номером или ошибку.
	Create(ctx context.Context, input *model.NewUser) (*model.User, error)

	// GetAll возвращает слайс всех пользователей или ошибку.
	// Если пользователей нет, то возвращается ошибка [errs.Empty].
	GetAll(ctx context.Context) ([]model.User, error)

	// Get возвращает пользователя по номеру или ошибку.
	// Если пользователь с таким номером не нашелся, то возвращается ошибка [errs.NotFound].
	Get(ctx context.Context, ID int) (*model.User, error)

	// Update обновляет пользователя и его данные для входа по номеру и возвращает его или ошибку.
	// Если пользователь с таким номером не нашелся, то возращается ошибка [errs.NotFound].
	Update(ctx context.Context, ID int, update *model.NewUser) (*model.User, error)

	// Delete удаляет пользователя по номеру и возвращает ошибку, если удаления не произошло.
	// Если пользователь с таким номером не нашелся, то возращается ошибка [errs.NotFound].
	Delete(ctx context.Context, ID int) error
}

// CreateUser получает данные о пользователе и его данных для входа
// из тела запроса и создает их. В ответе возвращается номер нового пользователя.
func (h *Handler) CreateUser(c echo.Context) error {
	newUser := new(model.NewUser)
	if err := bindAndValidate(c, newUser); err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("bind newUser: %w", err))
	}
	user, err := h.User.Create(c.Request().Context(), newUser)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"ID": user.ID,
	})
}

// GetAllUsers возвращает в ответе всех пользователей.
func (h *Handler) GetAllUsers(c echo.Context) error {
	users, err := h.User.GetAll(c.Request().Context())
	if err != nil && !errors.Is(err, errs.Empty) {
		return err
	}
	return c.JSON(http.StatusOK, users)
}

// GetUser получает номер пользователя из параметра id
// и в ответе возвращает пользователя с данным номером.
// Если пользователь с таким номером не нашелся,
// то возвращается ошибка [errs.NotFound].
func (h *Handler) GetUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("bind ID: %w", err))
	}
	user, err := h.User.Get(c.Request().Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

// UpdateUser получает номер пользователя из параметра id, данные о пользователя из тела запроса
// и обновляет пользователя с данным номером. В ответе ничего не возвращает.
func (h *Handler) UpdateUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("bind ID: %w", err))
	}
	userUpdate := new(model.NewUser)
	if err := bindAndValidate(c, userUpdate); err != nil {
		return echo.ErrBadRequest.WithInternal(err)
	}

	_, err = h.User.Update(c.Request().Context(), userID, userUpdate)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteUser получает номер пользователя из параметра id
// и удаляет пользователя с данным номером. В ответе ничего не возвращает.
func (h *Handler) DeleteUser(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("bind ID: %w", err))
	}
	err = h.User.Delete(c.Request().Context(), userID)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
