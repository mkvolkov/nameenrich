package aserver

import (
	"nameenrich/storage"

	"github.com/jmoiron/sqlx"
	"github.com/savsgio/atreugo/v11"
)

const (
	ResponseOK     = 200
	ResponseBadReq = 400
)

type Routes interface {
	GetMen() atreugo.View
}

type RBase struct {
	Psql *sqlx.DB
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
