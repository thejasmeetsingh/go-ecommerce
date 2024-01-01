// Common config modules

package api

import (
	"database/sql"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/thejasmeetsingh/go-ecommerce/order_service/internal/database"
)

type APIConfig struct {
	DB      *sql.DB
	Queries *database.Queries
	Cache   *redis.Client
}

func GetDBConn(dbURL string) *sql.DB {
	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		panic(fmt.Sprintf("Cannot connect to the database: %v", err))
	}

	return conn
}

func GetRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "orders_cache:6379",
		Password: "",
		DB:       0,
	})

	return client
}
