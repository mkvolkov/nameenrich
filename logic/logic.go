package logic

import (
	"nameenrich/enrich"
	"nameenrich/storage"
	"nameenrich/types"
	"strconv"

	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
)

type NewUser struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

func AddNewUser(user *NewUser, db *sqlx.DB, rconn redis.Conn) error {
	var baseMsg types.MsgBase
	baseMsg.Name = user.Name
	baseMsg.Surname = user.Surname
	baseMsg.Patronymic = user.Patronymic

	var enrMsg types.MsgEnriched = types.MsgEnriched{}

	enrich.Enrichment(&baseMsg, &enrMsg)

	err := storage.WriteData(db, rconn, &enrMsg)
	if err != nil {
		return err
	}

	return nil
}

func GetPersonByID(id int, db *sqlx.DB, rconn redis.Conn) (result []storage.GetPeopleResult, err error) {
	values, err := redis.Strings(rconn.Do("HGETALL", id))
	if err != nil {
		result, err = storage.GetPeopleByID(db, id)
		if err != nil {
			return result, err
		}

		return result, nil
	}
	if len(values) == 0 {
		result, err = storage.GetPeopleByID(db, id)
		if err != nil {
			return result, err
		} else {
			_, err = rconn.Do(
				"HMSET",
				result[0].ID,
				"p_name",
				result[0].Name,
				"surname",
				result[0].Surname,
				"patronymic",
				result[0].Patronymic,
				"age",
				result[0].Age,
				"gender",
				result[0].Gender,
			)
			if err != nil {
				return result, err
			}
		}

		return result, nil
	} else {
		res := storage.GetPeopleResult{}
		res.ID = id
		res.Name = values[1]
		res.Surname = values[3]
		res.Patronymic = values[5]
		age, err := strconv.Atoi(values[7])
		if err != nil {
			return result, err
		}
		res.Age = age
		res.Gender = values[9]
		result = append(result, res)
	}

	return result, nil
}

func GetCountryByName(name string, db *sqlx.DB, rconn redis.Conn) (result []storage.GetCountryResult, err error) {
	values, err := redis.Strings(rconn.Do("HGETALL", name))
	if err != nil {
		result, err = storage.GetCountryByName(db, name)
		if err != nil {
			return result, err
		}
		return result, err
	}

	for i := 0; i < len(values); i += 2 {
		countryRes := storage.GetCountryResult{}
		countryRes.CountryID = values[i]
		fVal, err := strconv.ParseFloat(values[i+1], 64)
		if err != nil {
			return result, err
		}
		countryRes.Probability = fVal

		result = append(result, countryRes)
	}

	return result, nil
}

func DeleteUser(id int, db *sqlx.DB, rconn redis.Conn) error {
	_, err := rconn.Do("DEL", id)
	if err != nil {
		return err
	}

	err = storage.DeleteUserByID(db, id)
	if err != nil {
		return err
	}

	return nil
}

func ChangeSurname(id int, surname string, db *sqlx.DB, rconn redis.Conn) error {
	_, err := rconn.Do("HSET", id, "surname", surname)
	if err != nil {
		return err
	}

	err = storage.ChangeSurname(db, id, surname)
	if err != nil {
		return err
	}

	return nil
}

func ChangeAge(id, age int, db *sqlx.DB, rconn redis.Conn) error {
	_, err := rconn.Do("HSET", id, "age", age)
	if err != nil {
		return err
	}

	err = storage.ChangeAge(db, id, age)
	if err != nil {
		return err
	}

	return nil
}
