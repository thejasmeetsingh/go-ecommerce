package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/thejasmeetsingh/go-ecommerce/order_service/api"
	"github.com/thejasmeetsingh/go-ecommerce/order_service/internal/database"
)

func main() {
	godotenv.Load()
	engine := gin.Default()

	port := os.Getenv("PORT")
	mode := os.Getenv("GIN_MODE")
	dbURL := os.Getenv("DB_URL")

	if port == "" {
		panic("Port is not configured")
	}

	if mode == "" || mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	if dbURL == "" {
		panic("DB URL is not configured")
	}

	dbConn := api.GetDBConn(dbURL)
	defer dbConn.Close()

	redisClient := api.GetRedisClient()
	defer redisClient.Close()

	apiCfg := api.APIConfig{
		DB:      dbConn,
		Queries: database.New(dbConn),
		Cache:   redisClient,
	}

	api.GetRoutes(engine, &apiCfg)

	log.Infoln("Order services started!")
	engine.Run(":" + port)
}
