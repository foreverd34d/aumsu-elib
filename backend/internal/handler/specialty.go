package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/labstack/echo/v4"
)

// SpecialtyService декларирует методы для работы со специальностями.
type SpecialtyService interface {
	// Create создает новую специальность и возвращает ее с номером или ошибку.
	Create(ctx context.Context, input *model.NewSpecialty) (*model.Specialty, error)

	// GetAll возвращает слайс всех специальностей или ошибку.
	// Если специальностей нет, то возвращается ошибка [errs.Empty].
	GetAll(ctx context.Context) ([]model.Specialty, error)

	// Get возвращает специальность по номеру или ошибку.
	// Если специальность с таким номером не нашлась, то возвращается ошибка [errs.NotFound].
	Get(ctx context.Context, ID int) (*model.Specialty, error)

	// Update обновляет специальность по номеру и возвращает ее с номером или ошибку.
	// Если специальность с таким номером не нашлась, то возращается ошибка [errs.NotFound].
	Update(ctx context.Context, ID int, update *model.NewSpecialty) (*model.Specialty, error)

	// Delete удаляет специальность по номеру и возвращает ошибку, если удаления не произошло.
	// Если специальность с таким номером не нашлась, то возращается ошибка [errs.NotFound].
	Delete(ctx context.Context, ID int) error
}

// CreateSpecialty получает данные о специальности из тела запроса и создает ее.
// В ответе возвращается номер новой кафедры.
func (h *Handler) CreateSpecialty(c echo.Context) error {
	newSpecialty := new(model.NewSpecialty)
	if err := bindAndValidate(c, newSpecialty); err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("bind newSpecialty: %w", err))
	}
	specialty, err := h.Specialty.Create(c.Request().Context(), newSpecialty)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"ID": specialty.ID,
	})
}

// GetAllSpecialties возвращает в ответе все специальности.
func (h *Handler) GetAllSpecialties(c echo.Context) error {
	specialties, err := h.Specialty.GetAll(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, specialties)
}

// GetSpecialty получает номер специальности из параметра id
// и в ответе возвращает специальность с данным номером.
// Если специальность с таким номером не нашлась,
// то возвращается ошибка [errs.NotFound].
func (h *Handler) GetSpecialty(c echo.Context) error {
	specialtyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("parse specialtyID: %w", err))
	}
	specialty, err := h.Specialty.Get(c.Request().Context(), specialtyID)
	if err != nil {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, specialty)
}

// UpdateSpecialty получает номер специальности из параметра id, данные о специальности из тела запроса
// и обновляет специальность с данным номером. В ответе ничего не возвращает.
func (h *Handler) UpdateSpecialty(c echo.Context) error {
	specialtyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("parse specialtyID: %w", err))
	}
	specialtyUpdate := new(model.NewSpecialty)
	if err := bindAndValidate(c, specialtyUpdate); err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("bind specialtyUpdate: %w", err))
	}
	_, err = h.Specialty.Update(c.Request().Context(), specialtyID, specialtyUpdate)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// DeleteSpecialty получает номер специальности из параметра id
// и удаляет специальность с данным номером. В ответе ничего не возвращает.
func (h *Handler) DeleteSpecialty(c echo.Context) error {
	specialtyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("parse specialtyID: %w", err))
	}
	if err := h.Specialty.Delete(c.Request().Context(), specialtyID); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
