package model

type User struct {
	UUID               string `json:"uuid"`
	FirstName          string `json:"first_name"`
	MiddleName         string `json:"middle_name"`
	LastName           string `json:"last_name"`
	CurrentDivision    string `json:"current_division"`
	CurrentDepartement string `json:"current_department"`
	CurrentPosition    string `json:"current_position"`
	Age                int    `json:"age"`
	Gender             string `json:"gender"`
	JoinDate           string `json:"join_date"`
	Status             string `json:"status"`
}
