package model

// Discipline представляет учебный предмет (дисциплину).
type Discipline struct {
	ID          int    `json:"disciplineID" db:"discipline_id"` // номер
	Name        string `json:"name" db:"name"`                  // название предмета
	SpecialtyID int    `json:"specialtyID" db:"specialty_id"`   // номер специальности
}
