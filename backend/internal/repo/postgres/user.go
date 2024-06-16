package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/foreverd34d/aumsu-elib/internal/errs"
	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/jmoiron/sqlx"
)

// UserRepo предоставляет доступ к базе данных с пользователями и их данными для входа.
type UserRepo struct {
	db *sqlx.DB
}

// NewUserRepo создает новый экземпляр [UserRepo].
func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create сохраняет нового пользователя и его данные для входа и возвращает пользователя с номером или ошибку.
func (ur *UserRepo) Create(ctx context.Context, input *model.NewUser) (*model.User, error) {
	txCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	tx, err := ur.db.BeginTxx(txCtx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w: %w", errs.Internal, err)
	}

	user := new(model.User)
	userQuery := `
		INSERT INTO users (surname, name, patronymic, role_id, group_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING user_id, surname, name, patronymic, role_id, group_id
	`
	if err := tx.GetContext(ctx, user, userQuery, input.Surname, input.Name, input.Patronymic, input.RoleID, input.GroupID); err != nil {
		return nil, fmt.Errorf("INSERT user: %w: %w", errs.Internal, err)
	}

	credentials := new(model.UserCredentials)
	credentialsQuery := `
		INSERT INTO users_credentials (login, password_hash)
		VALUES ($1, $2)
		RETURNING user_credentials_id, login, password_hash
	`
	if err := tx.GetContext(ctx, credentials, credentialsQuery, input.Login, input.Password); err != nil {
		return nil, fmt.Errorf("INSERT user's credentials: %w: %w", errs.Internal, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit changes: %w: %w", errs.Internal, err)
	}
	return user, nil
}

// GetAll возвращает слайс всех пользователей или ошибку.
// Если база данных пуста, то возвращается ошибка [errs.Empty].
func (ur *UserRepo) GetAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	query := `SELECT * FROM users`
	if err := ur.db.SelectContext(ctx, &users, query); err != nil {
		baseErr := errs.Internal
		if errors.Is(err, sql.ErrNoRows) {
			baseErr = errs.Empty
		}
		return nil, fmt.Errorf("SELECT users: %w: %w", baseErr, err)
	}
	return users, nil
}

// GetByID возвращает пользователя по номеру или ошибку.
// Если пользователь с таким номером не нашелся, то возвращается ошибка [errs.NotFound].
func (ur *UserRepo) GetByID(ctx context.Context, ID int) (*model.User, error) {
	user := new(model.User)
	query := `SELECT * FROM users WHERE user_id = $1`
	if err := ur.db.GetContext(ctx, user, query, ID); err != nil {
		baseErr := errs.Internal
		if errors.Is(err, sql.ErrNoRows) {
			baseErr = errs.NotFound
		}
		return nil, fmt.Errorf("SELECT user: %w: %w", baseErr, err)
	}
	return user, nil
}

// GetCredentialsByLogin возвращает данные пользователя для входа по логину или ошибку.
// Если пользователь с таким логином не нашелся, то возвращается ошибка [errs.InvalidLogin].
func (ur *UserRepo) GetCredentialsByLogin(ctx context.Context, login string) (*model.UserCredentials, error) {
	credentials := new(model.UserCredentials)
	query := `SELECT * FROM users_credentials WHERE login = $1`
	if err := ur.db.GetContext(ctx, credentials, query, login); err != nil {
		baseErr := errs.Internal
		if errors.Is(err, sql.ErrNoRows) {
			baseErr = errs.InvalidLogin
		}
		return nil, fmt.Errorf("SELECT user's credentials: %w: %w", baseErr, err)
	}
	return credentials, nil
}

// GetRole возвращает название роли пользователя или ошибку.
// Если пользователь с таким номером не нашелся, то возвращается ошибка [errs.NotFound].
func (ur *UserRepo) GetRole(ctx context.Context, userID int) (string, error) {
	var roleName string
	query := `
		SELECT r.name
		FROM roles r
		JOIN users u USING(role_id)
		WHERE u.user_id = $1
	`
	if err := ur.db.GetContext(ctx, &roleName, query, userID); err != nil {
		baseErr := errs.Internal
		if errors.Is(err, sql.ErrNoRows) {
			baseErr = errs.NotFound
		}
		return "", fmt.Errorf("SELECT user's role: %w: %w", baseErr, err)
	}
	return roleName, nil
}

// Update обновляет пользователя и его данные для входа по номеру и возвращает его или ошибку.
// Если пользователь с таким номером не нашелся, то возращается ошибка [errs.NotFound].
func (ur *UserRepo) Update(ctx context.Context, ID int, update *model.NewUser) (*model.User, error) {
	txCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	tx, err := ur.db.BeginTxx(txCtx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w: %w", errs.Internal, err)
	}

	updatedUser := new(model.User)
	userQuery := `
		UPDATE users
		SET surname = $1,
			name = $2,
			patronymic = $3,
			role_id = $6,
			group_id = $7
		WHERE user_id = $8
		RETURNING surname, name, patronymic, login, password_hash, role_id, group_id
	`
	if err := tx.GetContext(ctx, updatedUser, userQuery,
		update.Surname, update.Name, update.Patronymic, update.Login, update.Password, update.RoleID, update.GroupID, ID); err != nil {
		return nil, fmt.Errorf("UPDATE user: %w: %w", errs.NotFound, err)
	}

	credentialsQuery := `
		UPDATE credentials
		SET login = $1,
			password_hash = $2
		WHERE user_id = $3
	`
	_, err = tx.ExecContext(ctx, credentialsQuery, update.Login, update.Password)
	if err != nil {
		return nil, fmt.Errorf("UPDATE user's credentials: %w: %w", errs.NotFound, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit changes: %w: %w", errs.Internal, err)
	}

	return updatedUser, nil
}

// Delete удаляет пользователя по номеру и возвращает ошибку, если удаления не произошло.
// Если пользователь с таким номером не нашелся, то возращается ошибка [errs.NotFound].
func (ur *UserRepo) Delete(ctx context.Context, ID int) error {
	query := `DELETE FROM users WHERE user_id = $1 CASCADE`
	_, err := ur.db.ExecContext(ctx, query, ID)
	if err != nil {
		return fmt.Errorf("DELETE the user: %w: %w", errs.NotFound, err)
	}
	return err
}
