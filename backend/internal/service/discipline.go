package service

import (
	"context"
	"fmt"

	"github.com/foreverd34d/aumsu-elib/internal/model"
)

// DisciplineRepo определяет методы хранилища предметов.
type DisciplineRepo interface {
	// Create сохраняет предмет в хранилище и возвращает его с номером.
	Create(ctx context.Context, input *model.NewDiscipline) (*model.Discipline, error)

	// GetAll возвращает слайс всех предметов.
	// Если хранилище пусто, то возвращается ошибка [errs.Empty].
	GetAll(ctx context.Context) ([]model.Discipline, error)

	// Get возвращает предмет по номеру.
	// Если предмета с таким номером не нашлось, то возвращается ошибка [errs.NotFound].
	Get(ctx context.Context, ID int) (*model.Discipline, error)

	// Update обновляет предмет по номеру и возвращает его с номером.
	// Если предмета с таким номером не нашлось, то возвращается ошибка [errs.NotFound].
	Update(ctx context.Context, ID int, update *model.NewDiscipline) (*model.Discipline, error)

	// Delete удаляет предмет по номеру и возвращает ошибку, если удаления не произошло.
	// Если предмета с таким номером не нашлось, то возвращается ошибка [errs.NotFound].
	Delete(ctx context.Context, ID int) error
}

// DisciplineService реализует методы для работы с предметами
// и реализует интерфейс [handler.DisciplineService].
type DisciplineService struct {
	repo DisciplineRepo
}

// NewDisciplineService возвращает новый экземпляр [DisciplineService].
func NewDisciplineService(repo DisciplineRepo) *DisciplineService {
	return &DisciplineService{repo}
}

// Create создает новый предмет и возвращает ее с номером.
func (ds *DisciplineService) Create(ctx context.Context, input *model.NewDiscipline) (*model.Discipline, error) {
	discipline, err := ds.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("create new discipline: %w", err)
	}
	return discipline, nil
}

// GetAll возвращает слайс всех предметов. Если предметов нет, то возвращается ошибка [errs.Empty].
func (ds *DisciplineService) GetAll(ctx context.Context) ([]model.Discipline, error) {
	disciplines, err := ds.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all disciplines: %w", err)
	}
	return disciplines, nil
}

// Get возвращает предмет по номеру. Если предмета с таким номером не нашлось, то возращается ошибка [errs.NotFound].
func (ds *DisciplineService) Get(ctx context.Context, ID int) (*model.Discipline, error) {
	discipline, err := ds.repo.Get(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("get discipline: %w", err)
	}
	return discipline, nil
}

// Update обновляет предмет по номеру и возвращает обновленный предмет с номером.
// Если предмет с таким номером не нашелся, то возвращается ошибка [errs.NotFound].
func (ds *DisciplineService) Update(ctx context.Context, ID int, update *model.NewDiscipline) (*model.Discipline, error) {
	discipline, err := ds.repo.Update(ctx, ID, update)
	if err != nil {
		return nil, fmt.Errorf("update discipline: %w", err)
	}
	return discipline, nil
}

// Delete удаляет предмет по номеру и возвращает ошибку, если удаления не произошло.
// Если предмета с таким номером не нашлось, то возращается ошибка [errs.NotFound].
func (ds *DisciplineService) Delete(ctx context.Context, ID int) error {
	if err := ds.repo.Delete(ctx, ID); err != nil {
		return fmt.Errorf("delete discipline: %w", err)
	}
	return nil
}
