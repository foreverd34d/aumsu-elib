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

type DepartmentPsqlRepo struct {
	db *sqlx.DB
}

func NewDepartmentPsqlRepo(db *sqlx.DB) *DepartmentPsqlRepo {
	return &DepartmentPsqlRepo{db}
}

func (dr *DepartmentPsqlRepo) Create(ctx context.Context, input *model.NewDepartment) (*model.Department, error) {
	department := new(model.Department)
	query := `
		INSERT INTO departments (name)
		VALUES ($1)
		RETURNING department_id, name 
	`
	if err := dr.db.GetContext(ctx, department, query, input.Name); err != nil {
		return nil, fmt.Errorf("INSERT department: %w: %w", errs.Internal, err)
	}
	return department, nil
}

func (dr *DepartmentPsqlRepo) GetAll(ctx context.Context) ([]model.Department, error) {
	var departments []model.Department
	query := `SELECT * FROM departments`
	if err := dr.db.SelectContext(ctx, &departments, query); err != nil {
		baseErr := errs.Internal
		if errors.Is(err, sql.ErrNoRows) {
			baseErr = errs.Empty
		}
		return nil, fmt.Errorf("SELECT departments: %w: %w", baseErr, err)
	}
	return departments, nil
}
func (dr *DepartmentPsqlRepo) Get(ctx context.Context, ID int) (*model.Department, error) {
	department := new(model.Department)
	query := `SELECT * FROM departments WHERE department_id = $1`
	if err := dr.db.GetContext(ctx, department, query, ID); err != nil {
		baseErr := errs.Internal
		if errors.Is(err, sql.ErrNoRows) {
			baseErr = errs.NotFound
		}
		return nil, fmt.Errorf("SELECT department: %w: %w", baseErr, err)
	}
	return department, nil
}
func (dr *DepartmentPsqlRepo) Update(ctx context.Context, ID int, update *model.NewDepartment) (*model.Department, error) {
	department := new(model.Department)
	query := `
		UPDATE departments
		SET name = $1
		WHERE department_id = $2
		RETURNING department_id, name
	`
	if err := dr.db.GetContext(ctx, department, query, update.Name, ID); err != nil {
		return nil, fmt.Errorf("UPDATE department: %w: %w", errs.Internal, err)
	}
	return department, nil
}
func (dr *DepartmentPsqlRepo) Delete(ctx context.Context, ID int) error {
	query := `DELETE FROM departments WHERE department_id = $1 CASCADE`
	_, err := dr.db.ExecContext(ctx, query, ID)
	if err != nil {
		return fmt.Errorf("DELETE department: %w: %w", errs.Internal, err)
	}
	return err
}
