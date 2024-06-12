package service

import (
	"context"

	"github.com/foreverd34d/aumsu-elib/internal/model"
)

type groupRepo interface {
	Create(ctx context.Context, input *model.NewGroup) (*model.Group, error)
	GetAll(ctx context.Context) ([]model.Group, error)
	Get(ctx context.Context, ID int) (*model.Group, error)
	Update(ctx context.Context, ID int, update *model.NewGroup) (*model.Group, error)
	Delete(ctx context.Context, ID int) error
}

type GroupService struct {
	repo groupRepo
}

func NewGroupService(repo groupRepo) *GroupService {
	return &GroupService{repo}
}

func (gs *GroupService) Create(ctx context.Context, input *model.NewGroup) (*model.Group, error) {
	return gs.repo.Create(ctx, input)
}

func (gs *GroupService) GetAll(ctx context.Context) ([]model.Group, error) {
	return gs.repo.GetAll(ctx)
}

func (gs *GroupService) Get(ctx context.Context, ID int) (*model.Group, error) {
	return gs.repo.Get(ctx, ID)
}

func (gs *GroupService) Update(ctx context.Context, ID int, update *model.NewGroup) (*model.Group, error) {
	return gs.repo.Update(ctx, ID, update)
}

func (gs *GroupService) Delete(ctx context.Context, ID int) error {
	return gs.repo.Delete(ctx, ID)
}
