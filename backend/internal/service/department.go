package service

import (
	"context"
	"fmt"

	"github.com/foreverd34d/aumsu-elib/internal/model"
)

// DepartmentRepo определяет методы хранилища кафедр.
type DepartmentRepo interface {
	// Create сохраняет новую кафедру в хранилище и возвращает кафедру с порядковым номером или ошибку.
	Create(ctx context.Context, input *model.NewDepartment) (*model.Department, error)

	// GetAll возвращает слайс всех кафедр или ошибку.
	// Если хранилище пусто, то возвращается ошибка [errs.Empty].
	GetAll(ctx context.Context) ([]model.Department, error)

	// Get возвращает кафедру по номеру или ошибку.
	// Если кафедра с таким номером не нашлась, то возвращается ошибка [errs.NotFound].
	Get(ctx context.Context, ID int) (*model.Department, error)

	// Update обновляет кафедру по номеру и возвращает обновленную кафедру с номером или ошибку.
	// Если кафедра с таким номером не нашлась, то возращается ошибка [errs.NotFound].
	Update(ctx context.Context, ID int, update *model.NewDepartment) (*model.Department, error)

	// Delete удаляет кафедру по номеру и возвращает ошибку, если удаления не произошло.
	// Если кафедра с таким номером не нашлась, то возращается ошибка [errs.NotFound].
	Delete(ctx context.Context, ID int) error
}

// DepartmentService реализует методы для работы с кафедрами
// и реализует интерфейс [handler.departmentService].
type DepartmentService struct {
	repo DepartmentRepo
}

// NewDepartmentService возвращает новый экземпляр [DepartmentService].
func NewDepartmentService(repo DepartmentRepo) *DepartmentService {
	return &DepartmentService{repo}
}

// Create создает новую кафедру и возвращает ее с порядковым номером или ошибку.
func (ds *DepartmentService) Create(ctx context.Context, input *model.NewDepartment) (*model.Department, error) {
	department, err := ds.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("create new department: %w", err)
	}
	return department, nil
}

// GetAll возвращает слайс всех кафедр или ошибку. Если кафедр нет, то возвращается ошибка [errs.Empty].
func (ds *DepartmentService) GetAll(ctx context.Context) ([]model.Department, error) {
	departments, err := ds.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get departments: %w", err)
	}
	return departments, nil
}

// Get возвращает кафедру по номеру или ошибку. Если кафедра с таким номером не нашлась, то возвращается ошибка [errs.NotFound].
func (ds *DepartmentService) Get(ctx context.Context, ID int) (*model.Department, error) {
	department, err := ds.repo.Get(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("get the department with ID %v: %w", ID, err)
	}
	return department, nil
}

// Update обновляет кафедру по номеру и возвращает обновленную кафедру с номером или ошибку.
// Если кафедра с таким номером не нашлась, то возращается ошибка [errs.NotFound]
func (ds *DepartmentService) Update(ctx context.Context, ID int, update *model.NewDepartment) (*model.Department, error) {
	department, err := ds.repo.Update(ctx, ID, update)
	if err != nil {
		return nil, fmt.Errorf("update the department with ID %v: %w", ID, err)
	}
	return department, err
}

// Delete удаляет кафедру по номеру и возвращает ошибку, если удаления не произошло.
// Если кафедра с таким номером не нашлась, то возращается ошибка [errs.NotFound]
func (ds *DepartmentService) Delete(ctx context.Context, ID int) error {
	err := ds.repo.Delete(ctx, ID)
	if err != nil {
		return fmt.Errorf("delete department with ID %v: %w", ID, err)
	}
	return nil
}
