package api

import (
	"database/sql"
	"fmt"
)

func GetDBConn(dbURL string) *sql.DB {
	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		panic(fmt.Sprintf("Cannot connect to the database: %v", err))
	}

	return conn
}
