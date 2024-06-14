package service

import (
	"context"
	"fmt"

	"github.com/foreverd34d/aumsu-elib/internal/model"
)


type userRepo interface {
	Create(ctx context.Context, input *model.NewUser) (*model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
	GetByID(ctx context.Context, ID int) (*model.User, error)
	GetCredentialsByLogin(ctx context.Context, login string) (*model.UserCredentials, error)
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
	user, err := us.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("repo: create user: %w", err)
	}
	return user, nil
}

func (us *UserService) GetAll(ctx context.Context) ([]model.User, error) {
	users, err := us.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo: get all users: %w", err)
	}
	return users, nil
}

func (us *UserService) Get(ctx context.Context, ID int) (*model.User, error) {
	user, err := us.repo.GetByID(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("repo: get user with ID %v: %w", ID, err)
	}
	return user, nil
}

func (us *UserService) Update(ctx context.Context, ID int, update *model.NewUser) (*model.User, error) {
	update.Password = hashPassword(update.Password)
	user, err := us.repo.Update(ctx, ID, update)
	if err != nil {
		return nil, fmt.Errorf("repo: update the user with ID %v: %w", ID, err)
	}
	return user, nil
}

func (us *UserService) Delete(ctx context.Context, ID int) error {
	err := us.repo.Delete(ctx, ID)
	if err != nil {
		return fmt.Errorf("repo: delete the user with ID %v: %w", ID, err)
	}
	return nil
}
