package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/thejasmeetsingh/go-ecommerce/src/user_service/api"
	"github.com/thejasmeetsingh/go-ecommerce/src/user_service/internal/database"
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

	conn := api.GetDBConn(dbURL)
	defer conn.Close()

	apiCfg := api.APIConfig{
		DB:      conn,
		Queries: database.New(conn),
	}

	api.GetRoutes(engine, &apiCfg)

	log.Infoln("User services started!")
	engine.Run(":" + port)
}
