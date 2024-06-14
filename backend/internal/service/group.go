package service

import (
	"context"
	"fmt"

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
	group, err := gs.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("create new group: %w", err)
	}
	return group, nil
}

func (gs *GroupService) GetAll(ctx context.Context) ([]model.Group, error) {
	groups, err := gs.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all groups: %w", err)
	}
	return groups, nil
}

func (gs *GroupService) Get(ctx context.Context, ID int) (*model.Group, error) {
	group, err := gs.repo.Get(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("get the group with ID %v: %w", ID, err)
	}
	return group, nil
}

func (gs *GroupService) Update(ctx context.Context, ID int, update *model.NewGroup) (*model.Group, error) {
	group, err := gs.repo.Update(ctx, ID, update)
	if err != nil {
		return nil, fmt.Errorf("update the group with ID %v: %w", ID, err)
	}
	return group, nil
}

func (gs *GroupService) Delete(ctx context.Context, ID int) error {
	if err := gs.repo.Delete(ctx, ID); err != nil {
		return fmt.Errorf("delete the group with ID %v: %w", ID, err)
	}
	return nil
}
