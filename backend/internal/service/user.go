package service

import (
	"context"
	"errors"

	"github.com/foreverd34d/aumsu-elib/internal/model"
)

var (
	ErrInvalidPassword = errors.New("invalid password")
)

type userRepo interface {
	Create(ctx context.Context, input *model.NewUser) (*model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
	GetByID(ctx context.Context, ID int) (*model.User, error)
	GetByLogin(ctx context.Context, login string) (*model.User, error)
	GetRole(ctx context.Context, ID int) (string, error)
	Update(ctx context.Context, ID int, update *model.NewUser) (*model.User, error)
	Delete(ctx context.Context, ID int) error
}

type UserService struct {
	repo userRepo
}

func NewUserService(repo userRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) Create(ctx context.Context, input *model.NewUser) (*model.User, error) {
	input.Password = hashPassword(input.Password)
	return us.repo.Create(ctx, input)
}

func (us *UserService) GetAll(ctx context.Context) ([]model.User, error) {
	return us.repo.GetAll(ctx)
}

func (us *UserService) Get(ctx context.Context, ID int) (*model.User, error) {
	return us.repo.GetByID(ctx, ID)
}

func (us *UserService) Update(ctx context.Context, ID int, update *model.NewUser) (*model.User, error) {
	update.Password = hashPassword(update.Password)
	return us.repo.Update(ctx, ID, update)
}

func (us *UserService) Delete(ctx context.Context, ID int) error {
	return us.repo.Delete(ctx, ID)
}
