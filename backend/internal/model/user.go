package model

type User struct {
	ID           int     `json:"userID" db:"user_id"`
	Name         string  `json:"name" db:"name"`
	Surname      string  `json:"surname" db:"surname"`
	Patronymic   *string `json:"patronymic,omitempty" db:"patronymic"`
	Login        string  `json:"login" db:"login"`
	PasswordHash string  `json:"passwordHash" db:"password_hash"`
	RoleID       int     `json:"roleID" db:"role_id"`
	GroupID      *int    `json:"groupID,omitempty" db:"group_id"`
}

type Role struct {
	ID   int    `json:"roleID" db:"role_id"`
	Name string `json:"name" db:"name"`
}

type Group struct {
	ID          int    `json:"groupID" db:"group_id"`
	Name        string `json:"name" db:"name"`
	SpecialtyID int    `json:"specialtyID" db:"specialty_id"`
}

type Specialty struct {
	ID           int    `json:"specialtyID" db:"specialty_id"`
	Name         string `json:"name" db:"name"`
	DepartmentID int    `json:"departmentID" db:"department_id"`
}

type Department struct {
	ID   int    `json:"departmentID" db:"department_id"`
	Name string `json:"name" db:"name"`
}
