package storage

import (
	"context"
	"fmt"
	"nameenrich/types"

	"github.com/jmoiron/sqlx"
)

const (
	queryInsertPerson = `INSERT INTO people
		(p_name, surname, patronymic, age, gender)
		VALUES ($1, $2, $3, $4, $5)`

	queryInsertFirstName = `INSERT INTO first_names
		(name) VALUES ($1)`

	queryInsertNation = `INSERT INTO nation
		(user_name, country_id, probability)
		VALUES ($1, $2, $3)`

	querySelectMen = `SELECT *
		FROM people WHERE gender LIKE 'male'`
)

type Nationality struct {
	Country     string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type GetPeopleResult struct {
	ID         int    `db:"id"`
	Name       string `db:"p_name"`
	Surname    string `db:"surname"`
	Patronymic string `db:"patronymic"`
	Age        int    `db:"age"`
	Gender     string `db:"gender"`
}

type GetCountryResult struct {
	CountryID   string  `db:"country_id"`
	Probability float64 `db:"probability"`
}

func WriteData(db *sqlx.DB, params *types.MsgEnriched) error {
	var err error
	_, err = db.ExecContext(
		context.Background(),
		queryInsertFirstName,
		params.Name,
	)

	if err == nil {
		for i := 0; i < len(params.Nationalites); i++ {
			_, err = db.ExecContext(
				context.Background(),
				queryInsertNation,
				params.Name,
				params.Nationalites[i].Country,
				params.Nationalites[i].Probability,
			)
			if err != nil {
				fmt.Println("Error inserting nationality: ", err.Error())
			}
		}
	}

	_, err = db.ExecContext(
		context.Background(),
		queryInsertPerson,
		params.Name,
		params.Surname,
		params.Patronymic,
		params.Age,
		params.Gender,
	)

	return err
}

func GetMen(db *sqlx.DB) (result []GetPeopleResult, err error) {
	err = db.SelectContext(
		context.Background(),
		&result,
		querySelectMen,
	)

	if err != nil {
		return result, err
	}

	return result, nil
}
