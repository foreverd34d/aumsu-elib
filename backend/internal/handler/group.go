package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/labstack/echo/v4"
)

// GroupService декларирует методы для работы с группами.
type GroupService interface {
	// Create создает новую группу и возвращает ее с номером или ошибку.
	Create(ctx context.Context, input *model.NewGroup) (*model.Group, error)

	// GetAll возвращает слайс всех групп или ошибку. Если групп нет, то возвращается ошибка [errs.Empty].
	GetAll(ctx context.Context) ([]model.Group, error)

	// Get возвращает группу по номеру или ошибку. Если группа с таким номером не нашлась, то возвращается ошибка [errs.NotFound].
	Get(ctx context.Context, ID int) (*model.Group, error)

	// Update обновляет группу по номеру и возвращает обновленную группу с номером или ошибку.
	// Если группа с таким номером не нашлась, то возращается ошибка [errs.NotFound].
	Update(ctx context.Context, ID int, update *model.NewGroup) (*model.Group, error)

	// Delete удаляет группу по номеру и возвращает ошибку, если удаления не произошло.
	// Если группа с таким номером не нашлась, то возращается ошибка [errs.NotFound].
	Delete(ctx context.Context, ID int) error
}

// CreateGroup получает данные о группе из тела запроса и создает ее.
// В ответе возвращается номер новой группы.
func (h *Handler) CreateGroup(c echo.Context) error {
	newGroup := new(model.NewGroup)
	if err := bindAndValidate(c, newGroup); err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("bind newGroup: %w", err))
	}
	group, err := h.Group.Create(c.Request().Context(), newGroup)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"ID": group.ID,
	})
}

// GetAllGroups возвращает в ответе все кафедры.
func (h *Handler) GetAllGroups(c echo.Context) error {
	groups, err := h.Group.GetAll(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, groups)
}

// GetGroup получает номер группы из параметра id
// и в ответе возвращает группу с данным номером.
func (h *Handler) GetGroup(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("parse groupID: %w", err))
	}
	group, err := h.Group.Get(c.Request().Context(), groupID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, group)
}

// UpdateGroup получает номер группы из параметра id, данные о группе из тела запроса
// и обновляет группу с данным номером. В ответе ничего не возвращает.
func (h *Handler) UpdateGroup(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("parse groupID: %w", err))
	}
	groupUpdate := new(model.NewGroup)
	if err := bindAndValidate(c, groupUpdate); err != nil {
		return echo.ErrBadRequest
	}
	if _, err := h.Group.Update(c.Request().Context(), groupID, groupUpdate); err != nil {
		return echo.ErrNotFound
	}
	return c.NoContent(http.StatusNoContent)
}

// DeleteGroup получает номер группы из параметра id
// и удаляет группу с данным номером. В ответе ничего не возвращает.
func (h *Handler) DeleteGroup(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("parse groupID: %w", err))
	}
	if err := h.Group.Delete(c.Request().Context(), groupID); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
