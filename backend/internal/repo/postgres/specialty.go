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

// SpecialtyRepo предоставляет доступ к базе данных со специальностями.
type SpecialtyRepo struct {
	db *sqlx.DB
}

// NewSpecialtyRepo создает новый экземпляр [SpecialtyRepo].
func NewSpecialtyRepo(db *sqlx.DB) *SpecialtyRepo {
	return &SpecialtyRepo{db}
}

// Create сохраняет новую специальность в базе данных и возвращает ее с номером или ошибку.
func (sr *SpecialtyRepo) Create(ctx context.Context, input *model.NewSpecialty) (*model.Specialty, error) {
	specialty := new(model.Specialty)
	query := `
		INSERT INTO specialties (name, department_id)
		VALUES ($1, $2)
		RETURNING specialty_id, name, department_id
	`
	if err := sr.db.GetContext(ctx, specialty, query, input.Name, input.DepartmentID); err != nil {
		return nil, fmt.Errorf("INSERT specialty: %w: %w", errs.Internal, err)
	}
	return specialty, nil
}


// GetAll возвращает слайс всех специальностей или ошибку.
// Если база данных пуста, то возвращается ошибка [errs.Empty].
func (sr *SpecialtyRepo) GetAll(ctx context.Context) ([]model.Specialty, error) {
	var specialties []model.Specialty
	query := `SELECT * FROM specialties`
	if err := sr.db.SelectContext(ctx, &specialties, query); err != nil {
		baseErr := errs.Internal
		if errors.Is(err, sql.ErrNoRows) {
			baseErr = errs.Empty
		}
		return nil, fmt.Errorf("SELECT specialties: %w: %w", baseErr, err)
	}
	return specialties, nil
}

// Get возвращает специальность по номеру или ошибку.
// Если специальность с таким номером не нашлась, то возвращается ошибка [errs.NotFound].
func (sr *SpecialtyRepo) Get(ctx context.Context, ID int) (*model.Specialty, error) {
	specialty := new(model.Specialty)
	query := `SELECT * FROM specialties WHERE specialty_id = $1`
	if err := sr.db.GetContext(ctx, specialty, query, ID); err != nil {
		baseErr := errs.Internal
		if errors.Is(err, sql.ErrNoRows) {
			baseErr = errs.NotFound
		}
		return nil, fmt.Errorf("SELECT specialty: %w: %w", baseErr, err)
	}
	return specialty, nil
}

// Update обновляет специальность по номеру и возвращает ее с номером или ошибку.
// Если специальность с таким номером не нашлась, то возращается ошибка [errs.NotFound].
func (sr *SpecialtyRepo) Update(ctx context.Context, ID int, update *model.NewSpecialty) (*model.Specialty, error) {
	specialty := new(model.Specialty)
	query := `
		UPDATE specialties
		SET name = $1,
			department_id = $2
		WHERE specialty_id = $3
		RETURNING specialty_id, name, department_id
	`
	if err := sr.db.GetContext(ctx, specialty, query, update.Name, update.DepartmentID, ID); err != nil {
		return nil, fmt.Errorf("UPDATE specialty: %w: %w", errs.Internal, err)
	}
	return specialty, nil
}

// Delete удаляет специальность по номеру и возвращает ошибку, если удаления не произошло.
// Если специальность с таким номером не нашлась, то возращается ошибка [errs.NotFound].
func (sr *SpecialtyRepo) Delete(ctx context.Context, ID int) error {
	query := `DELETE FROM specialties WHERE specialty_id = $1 CASCADE`
	_, err := sr.db.ExecContext(ctx, query, ID)
	if err != nil {
		return fmt.Errorf("DELETE specialty: %w: %w", errs.Internal, err)
	}
	return err
}
