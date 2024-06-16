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

// GroupRepo предоставляет доступ к базе данных с группами.
type GroupRepo struct {
	db *sqlx.DB
}

// NewGroupRepo создает новый экземпляр [GroupRepo].
func NewGroupRepo(db *sqlx.DB) *GroupRepo {
	return &GroupRepo{db}
}

// Create сохраняет новую группу в базе данных и возвращает группу с номером или ошибку.
func (gs *GroupRepo) Create(ctx context.Context, input *model.NewGroup) (*model.Group, error) {
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

// GetAll возвращает слайс всех групп или ошибку.
// Если база данных пуста, то возвращается ошибка [errs.Empty].
func (gs *GroupRepo) GetAll(ctx context.Context) ([]model.Group, error) {
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

// Get возвращает группу по номеру или ошибку.
// Если группа с таким номером не нашлась, то возвращается ошибка [errs.NotFound].
func (gs *GroupRepo) Get(ctx context.Context, ID int) (*model.Group, error) {
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

// Update обновляет группу по номеру и возвращает обновленную группу с номером или ошибку.
// Если группа с таким номером не нашлась, то возращается ошибка [errs.NotFound].
func (gs *GroupRepo) Update(ctx context.Context, ID int, update *model.NewGroup) (*model.Group, error) {
	group := new(model.Group)
	query := `
		UPDATE groups
		SET name = $1,
			specialty_id = $2
		WHERE group_id = $3
		RETURNING group_id, name, specialty_id
	`
	if err := gs.db.GetContext(ctx, group, query, update.Name, update.SpecialtyID, ID); err != nil {
		return nil, fmt.Errorf("UPDATE group: %w: %w", errs.NotFound, err)
	}
	return group, nil
}

// Delete удаляет группу по номеру и возвращает ошибку, если удаления не произошло.
// Если группа с таким номером не нашлась, то возращается ошибка [errs.NotFound].
func (gs *GroupRepo) Delete(ctx context.Context, ID int) error {
	query := `DELETE FROM groups WHERE group_id = $1 CASCADE`
	_, err := gs.db.ExecContext(ctx, query, ID)
	if err != nil {
		return fmt.Errorf("DELETE group: %w: %w", errs.NotFound, err)
	}
	return nil
}
