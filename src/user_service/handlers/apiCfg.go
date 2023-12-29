package handlers

import (
	"database/sql"

	"github.com/thejasmeetsingh/go-ecommerce/src/user_service/internal/database"
)

type ApiConfig struct {
	DB      *sql.DB
	Queries *database.Queries
}
