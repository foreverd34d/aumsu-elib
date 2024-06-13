package postgres

import (
	"context"

	"github.com/foreverd34d/aumsu-elib/internal/model"

	"github.com/jmoiron/sqlx"
)

type DepartmentPostgresRepo struct {
	db *sqlx.DB
}

func NewDepartmentPostgresRepo(db *sqlx.DB) *DepartmentPostgresRepo {
	return &DepartmentPostgresRepo{db}
}

func (dr *DepartmentPostgresRepo) Create(ctx context.Context, input *model.NewDepartment) (*model.Department, error) {
	department := new(model.Department)
	query := `
		INSERT INTO departments (name)
		VALUES ($1)
		RETURNING department_id, name 
	`
	err := dr.db.GetContext(ctx, department, query, input.Name)
	return department, err

}

func (dr *DepartmentPostgresRepo) GetAll(ctx context.Context) ([]model.Department, error) {
	var departments []model.Department
	query := `SELECT * FROM departments`
	err := dr.db.SelectContext(ctx, &departments, query)
	return departments, err
}
func (dr *DepartmentPostgresRepo) Get(ctx context.Context, ID int) (*model.Department, error) {
	department := new(model.Department)
	query := `SELECT * FROM departments WHERE department_id = $1`
	err := dr.db.GetContext(ctx, department, query, ID)
	return department, err
}
func (dr *DepartmentPostgresRepo) Update(ctx context.Context, ID int, update *model.NewDepartment) (*model.Department, error) {
	department := new(model.Department)
	query := `
		UPDATE departments
		SET name = $1
		WHERE department_id = $2
		RETURNING department_id, name
	`
	err := dr.db.GetContext(ctx, department, query, update.Name, ID)
	return department, err
}
func (dr *DepartmentPostgresRepo) Delete(ctx context.Context, ID int) error {
	query := `DELETE FROM departments WHERE department_id = $1 CASCADE`
	_, err := dr.db.ExecContext(ctx, query, ID)
	return err
}
