package api

import (
	"database/sql"

	"github.com/thejasmeetsingh/go-ecommerce/src/user_service/internal/database"
)

type APIConfig struct {
	DB      *sql.DB
	Queries *database.Queries
}
