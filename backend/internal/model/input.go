package model

type NewUser struct {
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Patronymic *string `json:"patronymic,omitempty"`
	Login      string  `json:"login"`
	Password   string  `json:"password"`
	RoleID     int     `json:"roleID"`
	GroupID    *int    `json:"groupID,omitempty"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewToken struct {
	RefreshToken string `json:"refreshToken"`
	ExpiresAt    int    `json:"expiresAt"`
}

type NewGroup struct {
	Name        string `json:"name"`
	SpecialtyID int    `json:"specialtyID"`
}

type NewSpecialty struct {
	Name         string `json:"name"`
	DepartmentID int    `json:"departmentID"`
}

type NewDepartment struct {
	Name string `json:"name"`
}
