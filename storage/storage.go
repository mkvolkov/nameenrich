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

	querySelectWomen = `SELECT id, p_name, surname, patronymic, age
        FROM people WHERE gender LIKE 'female'`

	queryGetPeopleByName = `SELECT id, p_name, surname, patronymic, age, gender
        FROM people WHERE p_name ILIKE $1::text`

	queryGetPeopleByAgeLess = `SELECT id, p_name, surname, patronymic, age, gender
        FROM people WHERE age < $1
        ORDER BY (CASE WHEN $2 = true THEN age END) DESC,
                 (CASE WHEN $2 = false THEN age END) ASC`

	queryGetPeopleByAgeMore = `SELECT id, p_name, surname, patronymic, age, gender
        FROM people WHERE age > $1
        ORDER BY (CASE WHEN $2 = true THEN age END) DESC,
                 (CASE WHEN $2 = false THEN age END) ASC`

	queryGetCountryByName = `SELECT country_id, probability
        FROM nation WHERE user_name ILIKE $1::text`

	queryDeleteUserByID = `DELETE FROM people
        WHERE id = $1`

	queryChangeSurname = `UPDATE people
        SET surname = $2 WHERE id = $1`

	queryChangeAge = `UPDATE people
        SET age = $2 WHERE id = $1`
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

func GetWomen(db *sqlx.DB) (result []GetPeopleResult, err error) {
	err = db.SelectContext(
		context.Background(),
		&result,
		querySelectWomen,
	)

	if err != nil {
		return result, err
	}

	return result, nil
}

func GetPeopleByName(db *sqlx.DB, name string) (result []GetPeopleResult, err error) {
	err = db.SelectContext(
		context.Background(),
		&result,
		queryGetPeopleByName,
		name,
	)

	if err != nil {
		return result, err
	}

	return result, nil
}

func SelectPeopleByAge(db *sqlx.DB, age int, less, desc bool) (result []GetPeopleResult, err error) {
	if less {
		err = db.SelectContext(
			context.Background(),
			&result,
			queryGetPeopleByAgeLess,
			age,
			desc,
		)
	} else {
		err = db.SelectContext(
			context.Background(),
			&result,
			queryGetPeopleByAgeMore,
			age,
			desc,
		)
	}

	if err != nil {
		return result, err
	}

	return result, nil
}

func SelectCountryByName(db *sqlx.DB, name string) (result []GetCountryResult, err error) {
	err = db.SelectContext(
		context.Background(),
		&result,
		queryGetCountryByName,
		name,
	)

	if err != nil {
		return result, err
	}

	return result, nil
}

func DeleteUserByID(db *sqlx.DB, id int) error {
	_, err := db.ExecContext(
		context.Background(),
		queryDeleteUserByID,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func ChangeSurname(db *sqlx.DB, id int, surname string) error {
	_, err := db.ExecContext(
		context.Background(),
		queryChangeSurname,
		id,
		surname,
	)

	if err != nil {
		return err
	}

	return nil
}

func ChangeAge(db *sqlx.DB, id, age int) error {
	_, err := db.ExecContext(
		context.Background(),
		queryChangeAge,
		id,
		age,
	)

	if err != nil {
		return err
	}

	return nil
}
