package model

type NewUser struct {
	Name       string  `json:"name"`
	Surname    string  `json:"surname"`
	Patronymic *string `json:"patronymic,omitempty"`
	Login      string  `json:"login"`
	Password   string  `json:"password"`
	RoleName   string  `json:"roleName"`
	GroupName  *string `json:"groupName,omitempty"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewSession struct {
	RefreshToken string
	ExpiresAt int
	UserID int
}
