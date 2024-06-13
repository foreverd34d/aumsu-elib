package postgres

import (
	"context"

	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserPostgresRepo struct {
	db *sqlx.DB
}

func NewUserPostgresRepo(db *sqlx.DB) *UserPostgresRepo {
	return &UserPostgresRepo{db: db}
}

func (ur *UserPostgresRepo) Create(ctx context.Context, input *model.NewUser) (*model.User, error) {
	tx, err := ur.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	user := new(model.User)
	userQuery := `
		INSERT INTO users (surname, name, patronymic, role_id, group_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING user_id, surname, name, patronymic, role_id, group_id
	`
	if err := tx.GetContext(ctx, user, userQuery, input.Surname, input.Name, input.Patronymic, input.RoleID, input.GroupID); err != nil {
		tx.Rollback()
		return nil, err
	}

	credentials := new(model.UserCredentials)
	credentialsQuery := `
		INSERT INTO users_credentials (login, password_hash)
		VALUES ($1, $2)
		RETURNING user_credentials_id, login, password_hash
	`
	if err := tx.GetContext(ctx, credentials, credentialsQuery, input.Login, input.Password); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *UserPostgresRepo) GetAll(ctx context.Context) ([]model.User, error) {
	var users []model.User
	query := `SELECT * FROM users`
	err := ur.db.SelectContext(ctx, &users, query)
	return users, err
}

func (ur *UserPostgresRepo) GetByID(ctx context.Context, ID int) (*model.User, error) {
	user := new(model.User)
	query := `SELECT * FROM users WHERE user_id = $1`
	err := ur.db.GetContext(ctx, user, query, ID)
	return user, err
}

func (ur *UserPostgresRepo) GetCredentialsByLogin(ctx context.Context, login string) (*model.UserCredentials, error) {
	credentials := new(model.UserCredentials)
	query := `SELECT * FROM users_credentials WHERE login = $1`
	err := ur.db.GetContext(ctx, credentials, query, login)
	return credentials, err
}

func (ur *UserPostgresRepo) GetRole(ctx context.Context, userID int) (string, error) {
	var roleName string
	query := `
		SELECT r.name
		FROM roles r
		JOIN users u USING(role_id)
		WHERE u.user_id = $1
	`
	err := ur.db.GetContext(ctx, &roleName, query, userID)
	return roleName, err
}

func (ur *UserPostgresRepo) Update(ctx context.Context, ID int, update *model.NewUser) (*model.User, error) {
	updatedUser := new(model.User)
	query := `
		UPDATE users
		SET surname = $1,
			name = $2,
			patronymic = $3,
			login = $4,
			password_hash = $5,
			role_id = $6,
			group_id = $7
		WHERE user_id = $8
		RETURNING surname, name, patronymic, login, password_hash, role_id, group_id
	`
	err := ur.db.GetContext(ctx, updatedUser, query,
		update.Surname, update.Name, update.Patronymic, update.Login, update.Password, update.RoleID, update.GroupID, ID)
	return updatedUser, err
}

func (ur *UserPostgresRepo) Delete(ctx context.Context, ID int) error {
	query := `DELETE FROM users WHERE user_id = $1 CASCADE`
	_, err := ur.db.ExecContext(ctx, query, ID)
	return err
}
