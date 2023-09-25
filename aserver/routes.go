package aserver

import (
	"nameenrich/logic"
	"nameenrich/storage"

	"github.com/goccy/go-json"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/savsgio/atreugo/v11"
)

const (
	ResponseOK     = 200
	ResponseBadReq = 400
)

type Routes interface {
	GetMen() atreugo.View
	GetWomen() atreugo.View
	GetPeopleByID() atreugo.View
	GetPeopleByName() atreugo.View
	GetPeopleByAge() atreugo.View
	GetCountryByName() atreugo.View
	AddUser() atreugo.View
	DeleteUserByID() atreugo.View
	ChangeSurname() atreugo.View
	ChangeAge() atreugo.View
}

type RBase struct {
	Psql  *sqlx.DB
	Rconn redis.Conn
}

type FilterID struct {
	ID int `json:"id"`
}

type FilterName struct {
	Name string `json:"name"`
}

type FilterAge struct {
	Age  int  `json:"age"`
	Less bool `json:"less"`
	Desc bool `json:"desc"`
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

func (r *RBase) GetMen() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		data, err := storage.GetMen(r.Psql)
		if err != nil {
			return ctx.TextResponse("Error in GetMen\n", ResponseBadReq)
		}

		return ctx.JSONResponse(data, ResponseOK)
	}
}

func (r *RBase) GetWomen() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		data, err := storage.GetWomen(r.Psql)
		if err != nil {
			return ctx.TextResponse("Error in GetWomen", ResponseBadReq)
		}

		return ctx.JSONResponse(data, ResponseOK)
	}
}

func (r *RBase) GetPeopleByID() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		body := ctx.PostBody()

		filterID := &FilterID{}
		err := json.Unmarshal(body, filterID)
		if err != nil {
			return ctx.ErrorResponse(err, ResponseBadReq)
		}

		data, err := logic.GetPersonByID(filterID.ID, r.Psql, r.Rconn)
		if err != nil {
			return ctx.TextResponse("Error in GetPeopleByID", ResponseBadReq)
		}

		return ctx.JSONResponse(data, ResponseOK)
	}
}

func (r *RBase) GetPeopleByName() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		body := ctx.PostBody()

		filterName := &FilterName{}
		err := json.Unmarshal(body, filterName)
		if err != nil {
			return ctx.ErrorResponse(err, ResponseBadReq)
		}

		data, err := storage.GetPeopleByName(r.Psql, filterName.Name)
		if err != nil {
			return ctx.TextResponse("Error in GetPeopleByName", ResponseBadReq)
		}

		return ctx.JSONResponse(data, ResponseOK)
	}
}

func (r *RBase) GetPeopleByAge() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		body := ctx.PostBody()

		filterAge := &FilterAge{}
		err := json.Unmarshal(body, filterAge)
		if err != nil {
			return ctx.ErrorResponse(err, ResponseBadReq)
		}

		data, err := storage.GetPeopleByAge(r.Psql, filterAge.Age, filterAge.Less, filterAge.Desc)
		if err != nil {
			return ctx.TextResponse("Error in GetPeopleByAge", ResponseBadReq)
		}

		return ctx.JSONResponse(data, ResponseOK)
	}
}

func (r *RBase) GetCountryByName() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		body := ctx.PostBody()

		filterName := &FilterName{}
		err := json.Unmarshal(body, filterName)
		if err != nil {
			return ctx.ErrorResponse(err, ResponseBadReq)
		}

		//data, err := storage.GetCountryByName(r.Psql, filterName.Name)
		data, err := logic.GetCountryByName(filterName.Name, r.Psql, r.Rconn)
		if err != nil {
			return ctx.TextResponse("Error in GetCountryByName", ResponseBadReq)
		}

		return ctx.JSONResponse(data, ResponseOK)
	}
}

func (r *RBase) AddUser() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		body := ctx.PostBody()

		newUser := &logic.NewUser{}
		err := json.Unmarshal(body, newUser)
		if err != nil {
			return ctx.ErrorResponse(err, ResponseBadReq)
		}

		err = logic.AddNewUser(newUser, r.Psql, r.Rconn)
		if err != nil {
			return ctx.ErrorResponse(err, ResponseBadReq)
		}

		return ctx.TextResponse("User added successfully\n", ResponseOK)
	}
}

func (r *RBase) DeleteUserByID() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		body := ctx.PostBody()

		ID := &DeleteID{}
		err := json.Unmarshal(body, ID)
		if err != nil {
			return ctx.ErrorResponse(err, ResponseBadReq)
		}

		err = logic.DeleteUser(ID.ID, r.Psql, r.Rconn)
		if err != nil {
			return ctx.ErrorResponse(err, ResponseBadReq)
		}

		return ctx.TextResponse("User deleted successfully\n", ResponseOK)
	}
}

func (r *RBase) ChangeSurname() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		body := ctx.PostBody()

		SurnameData := &SetSurname{}
		err := json.Unmarshal(body, SurnameData)
		if err != nil {
			return ctx.ErrorResponse(err, ResponseBadReq)
		}

		err = logic.ChangeSurname(SurnameData.ID, SurnameData.Surname, r.Psql, r.Rconn)
		if err != nil {
			return ctx.ErrorResponse(err, ResponseBadReq)
		}

		return ctx.TextResponse("Surname changed successfully\n", ResponseOK)
	}
}

func (r *RBase) ChangeAge() atreugo.View {
	return func(ctx *atreugo.RequestCtx) error {
		body := ctx.PostBody()

		AgeStruct := &SetAge{}
		err := json.Unmarshal(body, AgeStruct)
		if err != nil {
			return ctx.ErrorResponse(err, ResponseBadReq)
		}

		err = logic.ChangeAge(AgeStruct.ID, AgeStruct.Age, r.Psql, r.Rconn)
		if err != nil {
			return ctx.ErrorResponse(err, ResponseBadReq)
		}

		return ctx.TextResponse("Age changed successfully\n", ResponseOK)
	}
}
