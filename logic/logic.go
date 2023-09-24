package logic

import (
	"nameenrich/enrich"
	"nameenrich/storage"
	"nameenrich/types"

	"github.com/jmoiron/sqlx"
)

type NewUser struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

func AddNewUser(user *NewUser, db *sqlx.DB) error {
	var baseMsg types.MsgBase
	baseMsg.Name = user.Name
	baseMsg.Surname = user.Surname
	baseMsg.Patronymic = user.Patronymic

	var enrMsg types.MsgEnriched = types.MsgEnriched{}

	enrich.Enrichment(&baseMsg, &enrMsg)

	err := storage.WriteData(db, &enrMsg)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(id int, db *sqlx.DB) error {
	err := storage.DeleteUserByID(db, id)
	if err != nil {
		return err
	}

	return nil
}

func ChangeSurname(id int, surname string, db *sqlx.DB) error {
	err := storage.ChangeSurname(db, id, surname)
	if err != nil {
		return err
	}

	return nil
}

func ChangeAge(id, age int, db *sqlx.DB) error {
	err := storage.ChangeAge(db, id, age)
	if err != nil {
		return err
	}

	return nil
}
