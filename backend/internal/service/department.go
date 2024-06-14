package service

import (
	"context"
	"fmt"

	"github.com/foreverd34d/aumsu-elib/internal/model"
)

type departmentRepo interface {
	Create(ctx context.Context, input *model.NewDepartment) (*model.Department, error)
	GetAll(ctx context.Context) ([]model.Department, error)
	Get(ctx context.Context, ID int) (*model.Department, error)
	Update(ctx context.Context, ID int, update *model.NewDepartment) (*model.Department, error)
	Delete(ctx context.Context, ID int) error
}

type DepartmentService struct {
	repo departmentRepo
}

func NewDepartmentService(repo departmentRepo) *DepartmentService {
	return &DepartmentService{repo}
}

func (ds *DepartmentService) Create(ctx context.Context, input *model.NewDepartment) (*model.Department, error) {
	department, err := ds.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("create new department: %w", err)
	}
	return department, nil
}

func (ds *DepartmentService) GetAll(ctx context.Context) ([]model.Department, error) {
	departments, err := ds.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get departments: %w", err)
	}
	return departments, nil
}

func (ds *DepartmentService) Get(ctx context.Context, ID int) (*model.Department, error) {
	department, err := ds.repo.Get(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("get the department with ID %v: %w", ID, err)
	}
	return department, nil
}

func (ds *DepartmentService) Update(ctx context.Context, ID int, update *model.NewDepartment) (*model.Department, error) {
	department, err := ds.repo.Update(ctx, ID, update)
	if err != nil {
		return nil, fmt.Errorf("update the department with ID %v: %w", ID, err)
	}
	return department, err
}

func (ds *DepartmentService) Delete(ctx context.Context, ID int) error {
	err := ds.repo.Delete(ctx, ID)
	if err != nil {
		return fmt.Errorf("delete department with ID %v: %w", ID, err)
	}
	return nil
}
