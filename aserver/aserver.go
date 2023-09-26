package aserver

import (
	"context"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/savsgio/atreugo/v11"
)

type AServer struct {
	Host      string
	Port      string
	AtrServer *atreugo.Atreugo
	Psql      *sqlx.DB
	Rconn     redis.Conn
}

func NewServer(host, port string, psql *sqlx.DB, rconn redis.Conn) *AServer {
	addr := host + ":" + port
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
		Rconn:     rconn,
	}
}

func (s *AServer) Run(ctx context.Context) error {
	aHandlers := &RBase{Psql: s.Psql, Rconn: s.Rconn}
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
	s.AtrServer.GET("/women", rs.GetWomen())
	s.AtrServer.GET("/id", rs.GetPeopleByID())
	s.AtrServer.GET("/name", rs.GetPeopleByName())
	s.AtrServer.GET("/age", rs.GetPeopleByAge())
	s.AtrServer.GET("/country", rs.GetCountryByName())
	s.AtrServer.POST("/adduser", rs.AddUser())
	s.AtrServer.DELETE("/deluser", rs.DeleteUserByID())
	s.AtrServer.POST("/chsurname", rs.ChangeSurname())
	s.AtrServer.POST("/chage", rs.ChangeAge())
}
