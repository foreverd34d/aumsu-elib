package postgres

import (
	"context"

	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/jmoiron/sqlx"
)

type GroupPostgresRepo struct {
	db *sqlx.DB
}

func NewGroupPostgresRepo(db *sqlx.DB) *GroupPostgresRepo {
	return &GroupPostgresRepo{db}
}

func (gs *GroupPostgresRepo) Create(ctx context.Context, input *model.NewGroup) (*model.Group, error) {
	group := new(model.Group)
	query := `
		INSERT INTO groups (name, specialty_id)
		VALUES ($1, $2)
		RETURNING group_id, name, specialty_id
	`
	err := gs.db.GetContext(ctx, group, query, input.Name, input.SpecialtyID)
	return group, err
}

func (gs *GroupPostgresRepo) GetAll(ctx context.Context) ([]model.Group, error) {
	var groups []model.Group
	query := `SELECT * FROM groups`
	err := gs.db.SelectContext(ctx, &groups, query)
	return groups, err
}

func (gs *GroupPostgresRepo) Get(ctx context.Context, ID int) (*model.Group, error) {
	group := new(model.Group)
	query := `SELECT * FROM groups WHERE group_id = $1`
	err := gs.db.GetContext(ctx, group, query, ID)
	return group, err
}

func (gs *GroupPostgresRepo) Update(ctx context.Context, ID int, update *model.NewGroup) (*model.Group, error) {
	group := new(model.Group)
	query := `
		UPDATE groups
		SET name = $1,
			specialty_id = $2
		WHERE group_id = $3
		RETURNING group_id, name, specialty_id
	`
	err := gs.db.GetContext(ctx, group, query, update.Name, update.SpecialtyID, ID)
	return group, err
}

func (gs *GroupPostgresRepo) Delete(ctx context.Context, ID int) error {
	query := `DELETE FROM groups WHERE group_id = $1 CASCADE`
	_, err := gs.db.ExecContext(ctx, query, ID)
	return err
}
