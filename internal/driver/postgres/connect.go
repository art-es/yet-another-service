package postgres

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

func Connect(url string) *sql.DB {
	var db *sql.DB
	var err error

	for i := 0; i < 30; i++ {
		if db, err = sql.Open("postgres", url); err == nil {
			return db
		}

		time.Sleep(time.Millisecond * 300)
	}

	panic(err)
}
