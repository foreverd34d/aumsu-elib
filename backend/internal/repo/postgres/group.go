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

type GroupPsqlRepo struct {
	db *sqlx.DB
}

func NewGroupPsqlRepo(db *sqlx.DB) *GroupPsqlRepo {
	return &GroupPsqlRepo{db}
}

func (gs *GroupPsqlRepo) Create(ctx context.Context, input *model.NewGroup) (*model.Group, error) {
	group := new(model.Group)
	query := `
		INSERT INTO groups (name, specialty_id)
		VALUES ($1, $2)
		RETURNING group_id, name, specialty_id
	`
	if err := gs.db.GetContext(ctx, group, query, input.Name, input.SpecialtyID); err != nil {
		return nil, fmt.Errorf("INSERT group: %w: %w", errs.Internal, err)
	}
	return group, nil
}

func (gs *GroupPsqlRepo) GetAll(ctx context.Context) ([]model.Group, error) {
	var groups []model.Group
	query := `SELECT * FROM groups`
	if err := gs.db.SelectContext(ctx, &groups, query); err != nil {
		baseErr := errs.Internal
		if errors.Is(err, sql.ErrNoRows) {
			baseErr = errs.Empty
		}
		return nil, fmt.Errorf("SELECT groups: %w: %w", baseErr, err)
	}
	return groups, nil
}

func (gs *GroupPsqlRepo) Get(ctx context.Context, ID int) (*model.Group, error) {
	group := new(model.Group)
	query := `SELECT * FROM groups WHERE group_id = $1`
	if err := gs.db.GetContext(ctx, group, query, ID); err != nil {
		baseErr := errs.Internal
		if errors.Is(err, sql.ErrNoRows) {
			baseErr = errs.NotFound
		}
		return nil, fmt.Errorf("SELECT group: %w: %w", baseErr, err)
	}
	return group, nil
}

func (gs *GroupPsqlRepo) Update(ctx context.Context, ID int, update *model.NewGroup) (*model.Group, error) {
	group := new(model.Group)
	query := `
		UPDATE groups
		SET name = $1,
			specialty_id = $2
		WHERE group_id = $3
		RETURNING group_id, name, specialty_id
	`
	if err := gs.db.GetContext(ctx, group, query, update.Name, update.SpecialtyID, ID); err != nil {
		return nil, fmt.Errorf("UPDATE group: %w: %w", errs.Internal, err)
	}
	return group, nil
}

func (gs *GroupPsqlRepo) Delete(ctx context.Context, ID int) error {
	query := `DELETE FROM groups WHERE group_id = $1 CASCADE`
	_, err := gs.db.ExecContext(ctx, query, ID)
	if err != nil {
		return fmt.Errorf("DELETE group: %w: %w", errs.Internal, err)
	}
	return nil
}
