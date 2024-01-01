package api

import (
	"database/sql"
	"fmt"

	"github.com/thejasmeetsingh/go-ecommerce/src/user_service/internal/database"
)

type APIConfig struct {
	DB      *sql.DB
	Queries *database.Queries
}

func GetDBConn(dbURL string) *sql.DB {
	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		panic(fmt.Sprintf("Cannot connect to the database: %v", err))
	}

	return conn
}
