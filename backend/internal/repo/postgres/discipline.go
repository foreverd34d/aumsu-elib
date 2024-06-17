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

// DisciplineRepo предоставляет доступ к базе данных с предметами.
type DisciplineRepo struct {
	db *sqlx.DB
}

// NewDisciplineRepo возвращает новый экземпляр [DisciplineRepo].
func NewDisciplineRepo(db *sqlx.DB) *DisciplineRepo {
	return &DisciplineRepo{db}
}

// Create сохраняет предмет в базе данных и возвращает его с номером.
func (dr *DisciplineRepo) Create(ctx context.Context, input *model.NewDiscipline) (*model.Discipline, error) {
	discipline := new(model.Discipline)
	query := `
		INSERT INTO disciplines (name, specialty_id)
		VALUES ($1, $2)
		RETURNING discipline_id, name, specialty_id
	`
	if err := dr.db.GetContext(ctx, discipline, query, input.Name, input.SpecialtyID); err != nil {
		return nil, fmt.Errorf("INSERT discipline: %w: %w", errs.Internal, err)
	}
	return discipline, nil
}

// GetAll возвращает слайс всех предметов.
// Если база данных пуста, то возвращается ошибка [errs.Empty].
func (dr *DisciplineRepo) GetAll(ctx context.Context) ([]model.Discipline, error) {
	var disciplines []model.Discipline
	query := `SELECT * FROM disciplines`
	if err := dr.db.SelectContext(ctx, &disciplines, query); err != nil {
		baseErr := errs.Internal
		if errors.Is(err, sql.ErrNoRows) {
			baseErr = errs.Empty
		}
		return nil, fmt.Errorf("SELECT disciplines: %w: %w", baseErr, err)
	}
	return disciplines, nil
}

// Get возвращает предмет по номеру.
// Если предмета с таким номером не нашлось, то возвращается ошибка [errs.NotFound].
func (dr *DisciplineRepo) Get(ctx context.Context, ID int) (*model.Discipline, error) {
	discipline := new(model.Discipline)
	query := `SELECT * FROM disciplines WHERE discipline_id = $1`
	if err := dr.db.GetContext(ctx, discipline, query, ID); err != nil {
		baseErr := errs.Internal
		if errors.Is(err, sql.ErrNoRows) {
			baseErr = errs.NotFound
		}
		return nil, fmt.Errorf("SELECT discipline: %w: %w", baseErr, err)
	}
	return discipline, nil
}

// Update обновляет предмет по номеру и возвращает его с номером.
// Если предмета с таким номером не нашлось, то возвращается ошибка [errs.NotFound].
func (dr *DisciplineRepo) Update(ctx context.Context, ID int, update *model.NewDiscipline) (*model.Discipline, error) {
	discipline := new(model.Discipline)
	query := `
		UPDATE disciplines
		SET name = $1,
			specialty_id = $2
		WHERE discipline_id = $3
		RETURNING discipline_id, name, specialty_id
	`
	if err := dr.db.GetContext(ctx, discipline, query, update.Name, update.SpecialtyID); err != nil {
		return nil, fmt.Errorf("UPDATE discipline: %w: %w", errs.NotFound, err)
	}
	return discipline, nil
}

// Delete удаляет предмет по номеру и возвращает ошибку, если удаления не произошло.
// Если предмета с таким номером не нашлось, то возвращается ошибка [errs.NotFound].
func (dr *DisciplineRepo) Delete(ctx context.Context, ID int) error {
	query := `DELETE FROM disciplines WHERE discipline_id = $1`
	if _, err := dr.db.ExecContext(ctx, query, ID); err != nil {
		return fmt.Errorf("DELETE discipline: %w: %w", errs.NotFound, err)
	}
	return nil
}
