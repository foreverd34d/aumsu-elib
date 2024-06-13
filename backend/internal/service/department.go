package service

import (
	"context"

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
	return ds.repo.Create(ctx, input)
}

func (ds *DepartmentService) GetAll(ctx context.Context) ([]model.Department, error) {
	return ds.repo.GetAll(ctx)
}

func (ds *DepartmentService) Get(ctx context.Context, ID int) (*model.Department, error) {
	return ds.repo.Get(ctx, ID)
}

func (ds *DepartmentService) Update(ctx context.Context, ID int, update *model.NewDepartment) (*model.Department, error) {
	return ds.repo.Update(ctx, ID, update)
}

func (ds *DepartmentService) Delete(ctx context.Context, ID int) error {
	return ds.repo.Delete(ctx, ID)
}
