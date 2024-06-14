package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/labstack/echo/v4"
)

type departmentService interface {
	Create(ctx context.Context, input *model.NewDepartment) (*model.Department, error)
	GetAll(ctx context.Context) ([]model.Department, error)
	Get(ctx context.Context, ID int) (*model.Department, error)
	Update(ctx context.Context, ID int, update *model.NewDepartment) (*model.Department, error)
	Delete(ctx context.Context, ID int) error
}

func (h *Handler) CreateDepartment(c echo.Context) error {
	newDepartment := new(model.NewDepartment)
	if err := c.Bind(newDepartment); err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("bind newDepartment: %w", err))
	}
	department, err := h.Department.Create(c.Request().Context(), newDepartment)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"ID": department.ID,
	})
}

func (h *Handler) GetAllDepartments(c echo.Context) error {
	departments, err := h.Department.GetAll(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, departments)
}

func (h *Handler) GetDepartment(c echo.Context) error {
	departmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("parse departmentID: %w", err))
	}
	department, err := h.Department.Get(c.Request().Context(), departmentID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, department)
}

func (h *Handler) UpdateDepartment(c echo.Context) error {
	departmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("parse departmentID: %w", err))
	}
	departmentUpdate := new(model.NewDepartment)
	if err := c.Bind(departmentUpdate); err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("bind departmentUpdate: %w", err))
	}
	_, err = h.Department.Update(c.Request().Context(), departmentID, departmentUpdate)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) DeleteDepartment(c echo.Context) error {
	departmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("parse departmentID: %w", err))
	}
	if err := h.Department.Delete(c.Request().Context(), departmentID); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
