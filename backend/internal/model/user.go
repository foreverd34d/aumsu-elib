package model

// User представляет пользователя библиотеки.
type User struct {
	ID         int     `json:"userID" db:"user_id"`                  // номер
	Name       string  `json:"name" db:"name"`                       // имя
	Surname    string  `json:"surname" db:"surname"`                 // фамилия
	Patronymic *string `json:"patronymic,omitempty" db:"patronymic"` // отчество (если имеется)
	RoleID     int     `json:"roleID" db:"role_id"`                  // номер роли
	GroupID    *int    `json:"groupID,omitempty" db:"group_id"`      // номер взвода (есть у студентов, отсутствует у остальных)
}

// UserCredentials представляет входные данные пользователя.
type UserCredentials struct {
	ID           int    `json:"userCredentialsID" db:"user_credential_id"` // номер
	Login        string `json:"login" db:"login"`                          // имя пользователя
	PasswordHash string `json:"passwordHash" db:"password_hash"`           // захэшированный алгоритмом sha-256 пароль
	UserID       int    `json:"userID" db:"user_id"`                       // номер пользователя
}

// Role представляет роль пользователя.
type Role struct {
	ID   int    `json:"roleID" db:"role_id"` // номер
	Name string `json:"name" db:"name"`      // название роли
}

// Group представляет взвод студентов.
type Group struct {
	ID          int    `json:"groupID" db:"group_id"`         // номер
	Name        string `json:"name" db:"name"`                // название группы
	SpecialtyID int    `json:"specialtyID" db:"specialty_id"` // номер специальности
}

// Specialty представляет специальность.
type Specialty struct {
	ID           int    `json:"specialtyID" db:"specialty_id"`   // номер
	Name         string `json:"name" db:"name"`                  // название специальности
	DepartmentID int    `json:"departmentID" db:"department_id"` // номер кафедры
}

// Department представляет кафедру.
type Department struct {
	ID   int    `json:"departmentID" db:"department_id"` // номер
	Name string `json:"name" db:"name"`                  // название
}
