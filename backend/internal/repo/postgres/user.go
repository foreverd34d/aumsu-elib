package postgres

import (
	"context"
	"libserver/internal/model"

	"github.com/jmoiron/sqlx"
)

type UserPostgresRepo struct {
	db *sqlx.DB
}

func NewUserPostgresRepo(db *sqlx.DB) *UserPostgresRepo {
	return &UserPostgresRepo{db: db}
}

func (ur *UserPostgresRepo) Create(ctx context.Context, input *model.NewUser) (*model.User, error) {
	user := new(model.User)
	query := `
		INSERT INTO users (surname, name, patronymic, login, password_hash, role_id, group_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING user_id, surname, name, patronymic, login, password_hash, role_id, group_id
	`
	err := ur.db.GetContext(ctx, user, query,
		input.Surname, input.Name, input.Patronymic, input.Login, input.Password, input.RoleID, input.GroupID)
	return user, err
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

func (ur *UserPostgresRepo) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	user := new(model.User)
	query := `SELECT * FROM users WHERE login = $1`
	err := ur.db.GetContext(ctx, user, query, login)
	return user, err
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
		update.Surname, update.Name, update.Patronymic, update.Login, update.Password, update.RoleID, update.GroupID)
	return updatedUser, err
}

func (ur *UserPostgresRepo) Delete(ctx context.Context, ID int) error {
	query := `DELETE FROM users WHERE user_id = $1`
	_, err := ur.db.ExecContext(ctx, query, ID)
	return err
}

func (ur *UserPostgresRepo) getRoleID(ctx context.Context, roleName string) (int, error) {
	var ID int
	query := `SELECT role_id FROM roles WHERE name = $1`
	err := ur.db.GetContext(ctx, &ID, query, roleName)
	return ID, err
}

func (ur *UserPostgresRepo) getGroupID(ctx context.Context, groupName string) (int, error) {
	var ID int
	query := `SELECT group_id FROM groups WHERE name = $1`
	err := ur.db.GetContext(ctx, &ID, query, groupName)
	return ID, err
}
