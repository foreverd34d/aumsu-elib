package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/labstack/echo/v4"
)

// DepartmentService декларирует методы для работы с кафедрами.
type DepartmentService interface {
	// Create создает новую кафедру и возвращает ее с порядковым номером или ошибку.
	Create(ctx context.Context, input *model.NewDepartment) (*model.Department, error)

	// GetAll возвращает слайс всех кафедр или ошибку. Если кафедр нет, то возвращается ошибка [errs.Empty].
	GetAll(ctx context.Context) ([]model.Department, error)
	
	// Get возвращает кафедру по номеру или ошибку. Если кафедра с таким номером не нашлась, то возвращается ошибка [errs.NotFound].
	Get(ctx context.Context, ID int) (*model.Department, error)

	// Update обновляет кафедру по номеру и возвращает обновленную кафедру с номером или ошибку.
	// Если кафедра с таким номером не нашлась, то возращается ошибка [errs.NotFound]
	Update(ctx context.Context, ID int, update *model.NewDepartment) (*model.Department, error)

	// Delete удаляет кафедру по номеру и возвращает ошибку, если удаления не произошло.
	// Если кафедра с таким номером не нашлась, то возращается ошибка [errs.NotFound]
	Delete(ctx context.Context, ID int) error
}

// CreateDepartment получает данные о кафедре из тела запроса и создает ее.
// В ответе возвращается номер новой кафедры.
func (h *Handler) CreateDepartment(c echo.Context) error {
	newDepartment := new(model.NewDepartment)
	if err := bindAndValidate(c, newDepartment); err != nil {
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

// GetAllDepartments возвращает в ответе все кафедры.
func (h *Handler) GetAllDepartments(c echo.Context) error {
	departments, err := h.Department.GetAll(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, departments)
}

// GetDepartment получает номер кафедры из параметра id
// и в ответе возвращает кафедру с данным номером.
// Если кафедра с таким номером не нашлась,
// то возвращается ошибка [errs.NotFound].
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

// UpdateDepartment получает номер кафедры из параметра id, данные о кафедре из тела запроса
// и обновляет кафедру с данным номером. В ответе ничего не возвращает.
func (h *Handler) UpdateDepartment(c echo.Context) error {
	departmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("parse departmentID: %w", err))
	}
	departmentUpdate := new(model.NewDepartment)
	if err := bindAndValidate(c, departmentUpdate); err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("bind departmentUpdate: %w", err))
	}
	_, err = h.Department.Update(c.Request().Context(), departmentID, departmentUpdate)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// DeleteDepartment получает номер кафедры из параметра id
// и удаляет кафедру с данным номером. В ответе ничего не возвращает.
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
