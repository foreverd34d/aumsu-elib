package handler

type Handler struct {
	User       userService
	Session    sessionService
	Group      groupService
	Specialty  specialtyService
	Department departmentService
}
