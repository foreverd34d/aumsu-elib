package model

type NewUser struct {
	Name       string  `json:"name" validate:"required"`
	Surname    string  `json:"surname" validate:"required"`
	Patronymic *string `json:"patronymic,omitempty"`
	Login      string  `json:"login" validate:"required"`
	Password   string  `json:"password" validate:"required"`
	RoleID     int     `json:"roleID" validate:"required"`
	GroupID    *int    `json:"groupID,omitempty"`
}

type Credentials struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type NewToken struct {
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int    `json:"expiresAt"`
}

type NewGroup struct {
	Name        string `json:"name" validate:"required"`
	SpecialtyID int    `json:"specialtyID" validate:"required"`
}

type NewSpecialty struct {
	Name         string `json:"name" validate:"required"`
	DepartmentID int    `json:"departmentID" validate:"required"`
}

type NewDepartment struct {
	Name string `json:"name" validate:"required"`
}
