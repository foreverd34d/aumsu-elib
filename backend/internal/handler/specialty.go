package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/labstack/echo/v4"
)

type specialtyService interface {
	Create(ctx context.Context, input *model.NewSpecialty) (*model.Specialty, error)
	GetAll(ctx context.Context) ([]model.Specialty, error)
	Get(ctx context.Context, ID int) (*model.Specialty, error)
	Update(ctx context.Context, ID int, update *model.NewSpecialty) (*model.Specialty, error)
	Delete(ctx context.Context, ID int) error
}

func (h *Handler) CreateSpecialty(c echo.Context) error {
	newSpecialty := new(model.NewSpecialty)
	if err := c.Bind(newSpecialty); err != nil {
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

func (h *Handler) GetAllSpecialties(c echo.Context) error {
	specialties, err := h.Specialty.GetAll(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, specialties)
}

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

func (h *Handler) UpdateSpecialty(c echo.Context) error {
	specialtyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("parse specialtyID: %w", err))
	}
	specialtyUpdate := new(model.NewSpecialty)
	if err := c.Bind(specialtyUpdate); err != nil {
		return echo.ErrBadRequest.WithInternal(fmt.Errorf("bind specialtyUpdate: %w", err))
	}
	_, err = h.Specialty.Update(c.Request().Context(), specialtyID, specialtyUpdate)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

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
