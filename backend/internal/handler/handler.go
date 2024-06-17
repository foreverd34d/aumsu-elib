package handler

import "github.com/labstack/echo/v4"

// Handler определяет методы обработчиков маршрутов.
type Handler struct {
	User       UserService
	Session    SessionService
	Group      GroupService
	Specialty  SpecialtyService
	Department DepartmentService
	Discipline DisciplineService
}

// bindAndValidate биндит структуру из тела запроса и проверяет ее.
func bindAndValidate(c echo.Context, i any) error {
	if err := c.Bind(i); err != nil {
		return err
	}
	return c.Validate(i)
}
