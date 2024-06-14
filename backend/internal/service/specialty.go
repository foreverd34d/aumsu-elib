package service

import (
	"context"
	"fmt"

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
	specialty, err := ss.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("create new specialty: %w", err)
	}
	return specialty, nil
}

func (ss *SpecialtyService) GetAll(ctx context.Context) ([]model.Specialty, error) {
	specialties, err := ss.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all specialties: %w", err)
	}
	return specialties, nil
}

func (ss *SpecialtyService) Get(ctx context.Context, ID int) (*model.Specialty, error) {
	specialty, err := ss.repo.Get(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("get specialty with ID %v: %w", ID, err)
	}
	return specialty, nil
}

func (ss *SpecialtyService) Update(ctx context.Context, ID int, update *model.NewSpecialty) (*model.Specialty, error) {
	specialty, err := ss.repo.Update(ctx, ID, update)
	if err != nil {
		return nil, fmt.Errorf("update the specialty with ID %v: %w", ID, err)
	}
	return specialty, nil
}

func (ss *SpecialtyService) Delete(ctx context.Context, ID int) error {
	if err := ss.repo.Delete(ctx, ID); err != nil {
		return fmt.Errorf("delete the specialty with ID %v: %w", ID, err)
	}
	return nil
}
