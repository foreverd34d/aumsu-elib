package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/labstack/echo/v4"
)

// DisciplineService определяет методы для работы с предметами.
type DisciplineService interface {
	// Create создает новый предмет и возвращает ее с номером.
	Create(ctx context.Context, input *model.NewDiscipline) (*model.Discipline, error)
	
	// GetAll возвращает слайс всех предметов. Если предметов нет, то возвращается ошибка [errs.Empty].
	GetAll(ctx context.Context) ([]model.Discipline, error)

	// Get возвращает предмет по номеру. Если предмета с таким номером не нашлось, то возращается ошибка [errs.NotFound].
	Get(ctx context.Context, ID int) (*model.Discipline, error)

	// Update обновляет предмет по номеру и возвращает обновленный предмет с номером.
	// Если предмет с таким номером не нашелся, то возвращается ошибка [errs.NotFound].
	Update(ctx context.Context, ID int, update *model.NewDiscipline) (*model.Discipline, error)

	// Delete удаляет предмет по номеру и возвращает ошибку, если удаления не произошло.
	// Если предмета с таким номером не нашлось, то возращается ошибка [errs.NotFound].
	Delete(ctx context.Context, ID int) error
}

// CreateDiscipline получает данные о предмете из тела запроса и создает его.
// В ответе возвращается номер нового предмета.
func (h *Handler) CreateDiscipline(c echo.Context) error {
	newDiscipline := new(model.NewDiscipline)
	if err := bindAndValidate(c, newDiscipline); err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("bind newDiscipline: %w", err))
	}
	discipline, err := h.Discipline.Create(c.Request().Context(), newDiscipline)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"ID": discipline.ID,
	})
}

// GetAllDisciplines возвращает в ответе все предметы.
func (h *Handler) GetAllDisciplines(c echo.Context) error {
	disciplines, err := h.Discipline.GetAll(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, disciplines)
}

// GetDiscipline получает номер предмета из параметра id
// и возвращает в ответе предмет с данным номером.
func (h *Handler) GetDiscipline(c echo.Context) error {
	disciplineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("parse disciplineID: %w", err))
	}
	discipline, err := h.Discipline.Get(c.Request().Context(), disciplineID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, discipline)
}

// UpdateDiscipline получает номер предмета из параметра id, данные о предмете из тела запроса
// и обновляет предмет с данным номером. В ответе ничего не возвращает.
func (h *Handler) UpdateDiscipline(c echo.Context) error {
	disciplineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("parse disciplineID: %w", err))
	}
	disciplineUpdate := new(model.NewDiscipline)
	if err := bindAndValidate(c, disciplineUpdate); err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("bind newDiscipline: %w", err))
	}
	if _, err = h.Discipline.Update(c.Request().Context(), disciplineID, disciplineUpdate); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

// DeleteDiscipline получает номер предмета из параметра id и удаляет предмет с данным номером.
// В ответе ничего не возвращает.
func (h *Handler) DeleteDiscipline(c echo.Context) error {
	disciplineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("parse disciplineID: %w", err))
	}
	if err := h.Discipline.Delete(c.Request().Context(), disciplineID); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
