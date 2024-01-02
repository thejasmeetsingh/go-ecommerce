package api

import (
	"database/sql"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
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
