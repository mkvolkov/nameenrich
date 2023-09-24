package aserver

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/savsgio/atreugo/v11"
)

type AServer struct {
	Host      string
	Port      string
	AtrServer *atreugo.Atreugo
	Psql      *sqlx.DB
}

func NewServer(host, port string, psql *sqlx.DB) *AServer {
	addr := fmt.Sprintf("%s:%s", host, port)
	aCfg := atreugo.Config{
		Addr: addr,

		GracefulShutdown: true,
	}

	aSrv := atreugo.New(aCfg)

	return &AServer{
		Host:      host,
		Port:      port,
		AtrServer: aSrv,
		Psql:      psql,
	}
}

func (s *AServer) Run(ctx context.Context) error {
	aHandlers := &RBase{Psql: s.Psql}
	s.MapHandlers(aHandlers)

	go func() {
		if err := s.AtrServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed Atreugo listen and serve: %v\n", err)
		}
	}()

	<-ctx.Done()

	return nil
}

func (s *AServer) MapHandlers(rs Routes) {
	s.AtrServer.GET("/men", rs.GetMen())
}
