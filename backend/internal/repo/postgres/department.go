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

// DepartmentRepo предоставляет доступ к базе данных с кафедрами.
type DepartmentRepo struct {
	db *sqlx.DB
}

// NewDepartmentRepo создает новый экземпляр [DepartmentRepo].
func NewDepartmentRepo(db *sqlx.DB) *DepartmentRepo {
	return &DepartmentRepo{db}
}

// Create сохраняет новую кафедру в базе данных и возвращает кафедру с порядковым номером или ошибку.
func (dr *DepartmentRepo) Create(ctx context.Context, input *model.NewDepartment) (*model.Department, error) {
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

// GetAll возвращает слайс всех кафедр или ошибку.
// Если база данных пуста, то возвращается ошибка [errs.Empty].
func (dr *DepartmentRepo) GetAll(ctx context.Context) ([]model.Department, error) {
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

// Get возвращает кафедру по номеру или ошибку.
// Если кафедра с таким номером не нашлась, то возвращается ошибка [errs.NotFound].
func (dr *DepartmentRepo) Get(ctx context.Context, ID int) (*model.Department, error) {
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

// Update обновляет кафедру по номеру и возвращает обновленную кафедру с номером или ошибку.
// Если кафедра с таким номером не нашлась, то возращается ошибка [errs.NotFound].
func (dr *DepartmentRepo) Update(ctx context.Context, ID int, update *model.NewDepartment) (*model.Department, error) {
	department := new(model.Department)
	query := `
		UPDATE departments
		SET name = $1
		WHERE department_id = $2
		RETURNING department_id, name
	`
	if err := dr.db.GetContext(ctx, department, query, update.Name, ID); err != nil {
		return nil, fmt.Errorf("UPDATE department: %w: %w", errs.NotFound, err)
	}
	return department, nil
}

// Delete удаляет кафедру по номеру и возвращает ошибку, если удаления не произошло.
// Если кафедра с таким номером не нашлась, то возращается ошибка [errs.NotFound].
func (dr *DepartmentRepo) Delete(ctx context.Context, ID int) error {
	query := `DELETE FROM departments WHERE department_id = $1 CASCADE`
	_, err := dr.db.ExecContext(ctx, query, ID)
	if err != nil {
		return fmt.Errorf("DELETE department: %w: %w", errs.NotFound, err)
	}
	return err
}
