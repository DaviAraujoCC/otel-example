package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)
func OpenConnection() (*sql.DB, error) {
    conn, err := sql.Open("postgres", "host=localhost port=5432 user=gopher password=1122 dbname=foobar sslmode=disable")
    if err != nil {
        panic(err)
    }
    err = conn.Ping()
    if err != nil {
        panic(err)
    }
    return conn, err
}