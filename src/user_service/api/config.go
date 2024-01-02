package api

import (
	"database/sql"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"github.com/thejasmeetsingh/go-ecommerce/src/user_service/internal/database"
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
		Addr:     "users_cache:6379",
		Password: "",
		DB:       0,
	})

	return client
}

func GetPromRequestTotal() *prometheus.CounterVec {
	httpRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"handler", "method"},
	)

	return httpRequestsTotal
}

func GetPromRequestDuration() *prometheus.HistogramVec {
	httpRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"handler", "method"},
	)
	return httpRequestDuration
}
