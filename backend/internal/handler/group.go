package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/labstack/echo/v4"
)

type groupService interface {
	Create(ctx context.Context, input *model.NewGroup) (*model.Group, error)
	GetAll(ctx context.Context) ([]model.Group, error)
	Get(ctx context.Context, ID int) (*model.Group, error)
	Update(ctx context.Context, ID int, update *model.NewGroup) (*model.Group, error)
	Delete(ctx context.Context, ID int) error
}

func (h *Handler) CreateGroup(c echo.Context) error {
	newGroup := new(model.NewGroup)
	if err := c.Bind(newGroup); err != nil {
		return echo.ErrBadRequest
	}
	group, err := h.Group.Create(c.Request().Context(), newGroup)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"ID": group.ID,
	})
}

func (h *Handler) GetAllGroups(c echo.Context) error {
	groups, err := h.Group.GetAll(c.Request().Context())
	if err != nil {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, groups)
}

func (h *Handler) GetGroup(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	}
	group, err := h.Group.Get(c.Request().Context(), groupID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, group)
}

func (h *Handler) UpdateGroup(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	}
	groupUpdate := new(model.NewGroup)
	if err := c.Bind(groupUpdate); err != nil {
		return echo.ErrBadRequest
	}
	if _, err := h.Group.Update(c.Request().Context(), groupID, groupUpdate); err != nil {
		return echo.ErrNotFound
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DeleteGroup(c echo.Context) error {
	groupID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest
	}
	if err := h.Group.Delete(c.Request().Context(), groupID); err != nil {
		return echo.ErrNotFound
	}
	return c.NoContent(http.StatusNoContent)
}
