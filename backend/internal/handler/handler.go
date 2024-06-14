package handler

import "github.com/labstack/echo/v4"

type Handler struct {
	User       userService
	Session    sessionService
	Group      groupService
	Specialty  specialtyService
	Department departmentService
}

func bindAndValidate(c echo.Context, i any) error {
	if err := c.Bind(i); err != nil {
		return err
	}
	return c.Validate(i)
}
