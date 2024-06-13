package postgres

import (
	"context"

	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/jmoiron/sqlx"
)

type SpecialtyPostgresRepo struct {
	db *sqlx.DB
}

func NewSpecialtyPostgresRepo(db *sqlx.DB) *SpecialtyPostgresRepo {
	return &SpecialtyPostgresRepo{db}
}

func (sr *SpecialtyPostgresRepo) Create(ctx context.Context, input *model.NewSpecialty) (*model.Specialty, error) {
	specialty := new(model.Specialty)
	query := `
		INSERT INTO specialties (name, department_id)
		VALUES ($1, $2)
		RETURNING specialty_id, name, department_id
	`
	err := sr.db.GetContext(ctx, specialty, query, input.Name, input.DepartmentID)
	return specialty, err
}

func (sr *SpecialtyPostgresRepo) GetAll(ctx context.Context) ([]model.Specialty, error) {
	var specialties []model.Specialty
	query := `SELECT * FROM specialties`
	err := sr.db.SelectContext(ctx, &specialties, query)
	return specialties, err
}

func (sr *SpecialtyPostgresRepo) Get(ctx context.Context, ID int) (*model.Specialty, error) {
	specialty := new(model.Specialty)
	query := `SELECT * FROM specialties WHERE specialty_id = $1`
	err := sr.db.GetContext(ctx, specialty, query, ID)
	return specialty, err
}

func (sr *SpecialtyPostgresRepo) Update(ctx context.Context, ID int, update *model.NewSpecialty) (*model.Specialty, error) {
	specialty := new(model.Specialty)
	query := `
		UPDATE specialties
		SET name = $1,
			department_id = $2
		WHERE specialty_id = $3
		RETURNING specialty_id, name, department_id
	`
	err := sr.db.GetContext(ctx, specialty, query, update.Name, update.DepartmentID, ID)
	return specialty, err
}

func (sr *SpecialtyPostgresRepo) Delete(ctx context.Context, ID int) error {
	query := `DELETE FROM specialties WHERE specialty_id = $1 CASCADE`
	_, err := sr.db.ExecContext(ctx, query, ID)
	return err
}
