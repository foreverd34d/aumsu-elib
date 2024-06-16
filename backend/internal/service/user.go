package service

import (
	"context"
	"fmt"

	"github.com/foreverd34d/aumsu-elib/internal/model"
)


// UserRepo определяет методы хранилища пользователей и данными для их входа.
type UserRepo interface {
	// Create создает нового пользователя и его данные для входа и возвращает пользователя с номером или ошибку.
	Create(ctx context.Context, input *model.NewUser) (*model.User, error)

	// GetAll возвращает слайс всех пользователей или ошибку.
	// Если база данных пуста, то возвращается ошибка [errs.Empty].
	GetAll(ctx context.Context) ([]model.User, error)

	// GetByID возвращает пользователя по номеру или ошибку.
	// Если пользователь с таким номером не нашелся, то возвращается ошибка [errs.NotFound].
	GetByID(ctx context.Context, ID int) (*model.User, error)

	// GetCredentialsByLogin возвращает данные пользователя для входа по логину или ошибку.
	// Если пользователь с таким логином не нашелся, то возвращается ошибка [errs.InvalidLogin].
	GetCredentialsByLogin(ctx context.Context, login string) (*model.UserCredentials, error)

	// GetRole возвращает название роли пользователя или ошибку.
	// Если пользователь с таким номером не нашелся, то возвращается ошибка [errs.NotFound].
	GetRole(ctx context.Context, ID int) (string, error)

	// Update обновляет пользователя и его данные для входа по номеру и возвращает его или ошибку.
	// Если пользователь с таким номером не нашелся, то возращается ошибка [errs.NotFound].
	Update(ctx context.Context, ID int, update *model.NewUser) (*model.User, error)

	// Delete удаляет пользователя по номеру и возвращает ошибку, если удаления не произошло.
	// Если пользователь с таким номером не нашелся, то возращается ошибка [errs.NotFound].
	Delete(ctx context.Context, ID int) error
}

// UserService реализует методы для работы с пользователями и их данными для входа
// и реализует интерфейс [handler.UserService].
type UserService struct {
	repo UserRepo
}

// NewUserService возвращает новый экземпляр [UserService].
func NewUserService(repo UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

// Create создает нового пользователя и его данные для входа и возвращает пользователя с номером или ошибку.
func (us *UserService) Create(ctx context.Context, input *model.NewUser) (*model.User, error) {
	input.Password = hashPassword(input.Password)
	user, err := us.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("repo: create user: %w", err)
	}
	return user, nil
}

// GetAll возвращает слайс всех пользователей или ошибку.
// Если пользователей нет, то возвращается ошибка [errs.Empty].
func (us *UserService) GetAll(ctx context.Context) ([]model.User, error) {
	users, err := us.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo: get all users: %w", err)
	}
	return users, nil
}

// Get возвращает пользователя по номеру или ошибку.
// Если пользователь с таким номером не нашелся, то возвращается ошибка [errs.NotFound].
func (us *UserService) Get(ctx context.Context, ID int) (*model.User, error) {
	user, err := us.repo.GetByID(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("repo: get user with ID %v: %w", ID, err)
	}
	return user, nil
}

// Update обновляет пользователя и его данные для входа по номеру и возвращает его или ошибку.
// Если пользователь с таким номером не нашелся, то возращается ошибка [errs.NotFound].
func (us *UserService) Update(ctx context.Context, ID int, update *model.NewUser) (*model.User, error) {
	update.Password = hashPassword(update.Password)
	user, err := us.repo.Update(ctx, ID, update)
	if err != nil {
		return nil, fmt.Errorf("repo: update the user with ID %v: %w", ID, err)
	}
	return user, nil
}

// Delete удаляет пользователя по номеру и возвращает ошибку, если удаления не произошло.
// Если пользователь с таким номером не нашелся, то возращается ошибка [errs.NotFound].
func (us *UserService) Delete(ctx context.Context, ID int) error {
	err := us.repo.Delete(ctx, ID)
	if err != nil {
		return fmt.Errorf("repo: delete the user with ID %v: %w", ID, err)
	}
	return nil
}
