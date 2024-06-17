package model

// NewUser содержит данные для добавления нового пользователя.
type NewUser struct {
	Name       string  `json:"name" validate:"required"`           // имя
	Surname    string  `json:"surname" validate:"required"`        // фамилия
	Patronymic *string `json:"patronymic,omitempty"`               // отчество (если имеется)
	Login      string  `json:"login" validate:"required"`          // имя пользователя
	Password   string  `json:"password" validate:"required"`       // пароль
	RoleID     int     `json:"roleID" validate:"required,gte=1"`   // номер роли
	GroupID    *int    `json:"groupID,omitempty" validate:"gte=1"` // номер группы (должен быть у студента)
}

// Credentials содержит данные для входа в систему.
type Credentials struct {
	Username string `json:"username" validate:"required"` // имя пользователя
	Password string `json:"password" validate:"required"` // пароль
}

// NewToken содержит данные для создания нового токена обновления.
type NewToken struct {
	RefreshToken string `json:"refreshToken"` // токен обновления
	ExpiresAt    int    `json:"expiresAt"`    // время истечения срока действия токена, хранится в формате unix
}

// NewGroup содержит данные для добавления нового взвода.
type NewGroup struct {
	Name        string `json:"name" validate:"required"`              // название
	SpecialtyID int    `json:"specialtyID" validate:"required,gte=1"` // номер специальности
}

// NewSpecialty содержит данные для добавления новой специальности.
type NewSpecialty struct {
	Name         string `json:"name" validate:"required"`               // название
	DepartmentID int    `json:"departmentID" validate:"required,gte=1"` // номер кафедры
}

// NewDepartment содержит данные для добавления новой кафедры.
type NewDepartment struct {
	Name string `json:"name" validate:"required"` // название
}

type NewDiscipline struct {
	Name        string `json:"name" validate:"required"`
	SpecialtyID int    `json:"specialtyID" validate:"required,gte=1"`
}
