package types

type MsgBase struct {
	Name       string `json:"name" validate:"required"`
	Surname    string `json:"surname" validate:"required"`
	Patronymic string
}

type MsgError struct {
	Errname      string
	IncorrectMsg string
}

type MsgEnriched struct {
	Name         string
	Surname      string
	Patronymic   string
	Age          int
	Gender       string
	Nationalites []Nationality
}

type AgeData struct {
	Age int `json:"age"`
}

type GenderData struct {
	Gender string `json:"gender"`
}

type Nationality struct {
	Country     string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type NationData struct {
	Nationalities []Nationality `json:"country"`
}

type NewUser struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type DeleteID struct {
	ID int `json:"id"`
}

type SetSurname struct {
	ID      int    `json:"id"`
	Surname string `json:"surname"`
}

type SetAge struct {
	ID  int `json:"id"`
	Age int `json:"age"`
}
