package service

import (
	"context"

	"github.com/foreverd34d/aumsu-elib/internal/model"
)

type specialtyRepo interface {
	Create(ctx context.Context, input *model.NewSpecialty) (*model.Specialty, error)
	GetAll(ctx context.Context) ([]model.Specialty, error)
	Get(ctx context.Context, ID int) (*model.Specialty, error)
	Update(ctx context.Context, ID int, update *model.NewSpecialty) (*model.Specialty, error)
	Delete(ctx context.Context, ID int) error
}

type SpecialtyService struct {
	repo specialtyRepo
}

func NewSpecialtyService(repo specialtyRepo) *SpecialtyService {
	return &SpecialtyService{repo}
}

func (ss *SpecialtyService) Create(ctx context.Context, input *model.NewSpecialty) (*model.Specialty, error) {
	return ss.repo.Create(ctx, input)
}

func (ss *SpecialtyService) GetAll(ctx context.Context) ([]model.Specialty, error) {
	return ss.repo.GetAll(ctx)
}

func (ss *SpecialtyService) Get(ctx context.Context, ID int) (*model.Specialty, error) {
	return ss.repo.Get(ctx, ID)
}

func (ss *SpecialtyService) Update(ctx context.Context, ID int, update *model.NewSpecialty) (*model.Specialty, error) {
	return ss.repo.Update(ctx, ID, update)
}

func (ss *SpecialtyService) Delete(ctx context.Context, ID int) error {
	return ss.repo.Delete(ctx, ID)
}
