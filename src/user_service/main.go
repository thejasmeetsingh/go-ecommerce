package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/thejasmeetsingh/go-ecommerce/src/user_service/handlers"
	"github.com/thejasmeetsingh/go-ecommerce/src/user_service/internal/database"
	middlewares "github.com/thejasmeetsingh/go-ecommerce/src/user_service/middleware"
)

func loadRoutes(engine *gin.Engine, apiCfg *handlers.ApiConfig) {
	// Default route for health check
	engine.GET("/health-check/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Up and Running!",
		})
	})

	router := engine.Group("/api/v1/")
	authRouter := router.Group("")
	authRouter.Use(middlewares.JWTAuth(apiCfg))

	// Non auth routes
	router.POST("register/", apiCfg.Singup)
	router.POST("login/", apiCfg.Login)
	router.POST("refresh-token/", apiCfg.RefreshAccessToken)

	// Auth routes
	authRouter.GET("profile/", apiCfg.GetUserProfile)
	authRouter.PATCH("profile/", apiCfg.UpdateUserProfile)
	authRouter.DELETE("profile/", apiCfg.DeleteUserProfile)
	authRouter.PUT("change-password/", apiCfg.ChangePassword)
}

func getDBConn(dbURL string) *database.Queries {
	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		panic(fmt.Sprintf("Cannot connect to the database: %v", err))
	}

	return database.New(conn)
}

func main() {
	godotenv.Load()
	engine := gin.Default()

	port := os.Getenv("PORT")
	if port == "" {
		panic("Port is not configured")
	}

	mode := os.Getenv("GIN_MODE")

	if mode == "" || mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		panic("DB URL is not configured")
	}

	dbConn := getDBConn(dbURL)
	apiCfg := handlers.ApiConfig{
		DB: dbConn,
	}

	loadRoutes(engine, &apiCfg)

	log.Infoln("User services started!")
	engine.Run(":" + port)
}
