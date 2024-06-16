package service

import (
	"context"
	"fmt"

	"github.com/foreverd34d/aumsu-elib/internal/model"
)

// GroupRepo определяет методы хранилища взводов.
type GroupRepo interface {
	// Create сохраняет новую группу в хранилище и возвращает группу с номером или ошибку.
	Create(ctx context.Context, input *model.NewGroup) (*model.Group, error)

	// GetAll возвращает слайс всех групп или ошибку.
	// Если база данных пуста, то возвращается ошибка [errs.Empty].
	GetAll(ctx context.Context) ([]model.Group, error)

	// Get возвращает группу по номеру или ошибку.
	// Если группа с таким номером не нашлась, то возвращается ошибка [errs.NotFound].
	Get(ctx context.Context, ID int) (*model.Group, error)

	// Update обновляет группу по номеру и возвращает обновленную группу с номером или ошибку.
	// Если группа с таким номером не нашлась, то возращается ошибка [errs.NotFound].
	Update(ctx context.Context, ID int, update *model.NewGroup) (*model.Group, error)

	// Delete удаляет группу по номеру и возвращает ошибку, если удаления не произошло.
	// Если группа с таким номером не нашлась, то возращается ошибка [errs.NotFound].
	Delete(ctx context.Context, ID int) error
}

// GroupService определяет методы для работы с группами
// и реализует интерфейс [handler.groupService].
type GroupService struct {
	repo GroupRepo
}

// NewGroupService возвращает новый экземпляр [GroupService]
func NewGroupService(repo GroupRepo) *GroupService {
	return &GroupService{repo}
}

// Create создает новую группу и возвращает ее с номером или ошибку.
func (gs *GroupService) Create(ctx context.Context, input *model.NewGroup) (*model.Group, error) {
	group, err := gs.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("create new group: %w", err)
	}
	return group, nil
}

// GetAll возвращает слайс всех групп или ошибку. Если групп нет, то возвращается ошибка [errs.Empty].
func (gs *GroupService) GetAll(ctx context.Context) ([]model.Group, error) {
	groups, err := gs.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all groups: %w", err)
	}
	return groups, nil
}

// Get возвращает группу по номеру или ошибку. Если группа с таким номером не нашлась, то возвращается ошибка [errs.NotFound].
func (gs *GroupService) Get(ctx context.Context, ID int) (*model.Group, error) {
	group, err := gs.repo.Get(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("get the group with ID %v: %w", ID, err)
	}
	return group, nil
}

// Update обновляет группу по номеру и возвращает обновленную группу с номером или ошибку.
// Если группа с таким номером не нашлась, то возращается ошибка [errs.NotFound].
func (gs *GroupService) Update(ctx context.Context, ID int, update *model.NewGroup) (*model.Group, error) {
	group, err := gs.repo.Update(ctx, ID, update)
	if err != nil {
		return nil, fmt.Errorf("update the group with ID %v: %w", ID, err)
	}
	return group, nil
}

// Delete удаляет группу по номеру и возвращает ошибку, если удаления не произошло.
// Если группа с таким номером не нашлась, то возращается ошибка [errs.NotFound].
func (gs *GroupService) Delete(ctx context.Context, ID int) error {
	if err := gs.repo.Delete(ctx, ID); err != nil {
		return fmt.Errorf("delete the group with ID %v: %w", ID, err)
	}
	return nil
}
