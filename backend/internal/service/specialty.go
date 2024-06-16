package service

import (
	"context"
	"fmt"

	"github.com/foreverd34d/aumsu-elib/internal/model"
)

// SpecialtyRepo определяет методы хранилища специальностей.
type SpecialtyRepo interface {
	// Create сохраняет новую специальность в хранилище и возвращает ее с номером или ошибку.
	Create(ctx context.Context, input *model.NewSpecialty) (*model.Specialty, error)

	// GetAll возвращает слайс всех специальностей или ошибку.
	// Если база данных пуста, то возвращается ошибка [errs.Empty].
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

// SpecialtyService реализует методы для работы со специальностями
// и реализует интерфейс [handler.SpecialtyService].
type SpecialtyService struct {
	repo SpecialtyRepo
}

// NewSpecialtyService возвращает новый экземпляр [SpecialtyService].
func NewSpecialtyService(repo SpecialtyRepo) *SpecialtyService {
	return &SpecialtyService{repo}
}

// Create создает новую специальность и возвращает ее с номером или ошибку.
func (ss *SpecialtyService) Create(ctx context.Context, input *model.NewSpecialty) (*model.Specialty, error) {
	specialty, err := ss.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("create new specialty: %w", err)
	}
	return specialty, nil
}

// GetAll возвращает слайс всех специальностей или ошибку.
// Если специальностей нет, то возвращается ошибка [errs.Empty].
func (ss *SpecialtyService) GetAll(ctx context.Context) ([]model.Specialty, error) {
	specialties, err := ss.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all specialties: %w", err)
	}
	return specialties, nil
}

// Get возвращает специальность по номеру или ошибку.
// Если специальность с таким номером не нашлась, то возвращается ошибка [errs.NotFound].
func (ss *SpecialtyService) Get(ctx context.Context, ID int) (*model.Specialty, error) {
	specialty, err := ss.repo.Get(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("get specialty with ID %v: %w", ID, err)
	}
	return specialty, nil
}

// Update обновляет специальность по номеру и возвращает ее с номером или ошибку.
// Если специальность с таким номером не нашлась, то возращается ошибка [errs.NotFound].
func (ss *SpecialtyService) Update(ctx context.Context, ID int, update *model.NewSpecialty) (*model.Specialty, error) {
	specialty, err := ss.repo.Update(ctx, ID, update)
	if err != nil {
		return nil, fmt.Errorf("update the specialty with ID %v: %w", ID, err)
	}
	return specialty, nil
}

// Delete удаляет специальность по номеру и возвращает ошибку, если удаления не произошло.
// Если специальность с таким номером не нашлась, то возращается ошибка [errs.NotFound].
func (ss *SpecialtyService) Delete(ctx context.Context, ID int) error {
	if err := ss.repo.Delete(ctx, ID); err != nil {
		return fmt.Errorf("delete the specialty with ID %v: %w", ID, err)
	}
	return nil
}
