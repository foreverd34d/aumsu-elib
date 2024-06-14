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

type SpecialtyPsqlRepo struct {
	db *sqlx.DB
}

func NewSpecialtyPsqlRepo(db *sqlx.DB) *SpecialtyPsqlRepo {
	return &SpecialtyPsqlRepo{db}
}

func (sr *SpecialtyPsqlRepo) Create(ctx context.Context, input *model.NewSpecialty) (*model.Specialty, error) {
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

func (sr *SpecialtyPsqlRepo) GetAll(ctx context.Context) ([]model.Specialty, error) {
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

func (sr *SpecialtyPsqlRepo) Get(ctx context.Context, ID int) (*model.Specialty, error) {
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

func (sr *SpecialtyPsqlRepo) Update(ctx context.Context, ID int, update *model.NewSpecialty) (*model.Specialty, error) {
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

func (sr *SpecialtyPsqlRepo) Delete(ctx context.Context, ID int) error {
	query := `DELETE FROM specialties WHERE specialty_id = $1 CASCADE`
	_, err := sr.db.ExecContext(ctx, query, ID)
	if err != nil {
		return fmt.Errorf("DELETE specialty: %w: %w", errs.Internal, err)
	}
	return err
}
