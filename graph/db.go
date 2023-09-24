package graph

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(url string) *sqlx.DB {
	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil
	}

	err = db.Ping()
	if err != nil {
		return nil
	}

	return db
}
