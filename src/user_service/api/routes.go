package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRoutes(engine *gin.Engine, apiCfg *APIConfig) *gin.RouterGroup {
	// Default route for health check
	engine.GET("/health-check/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Up and Running!",
		})
	})

	router := engine.Group("/api/v1/")
	authRouter := router.Group("")
	authRouter.Use(JWTAuth((apiCfg)))

	// Non auth routes
	router.POST("register/", apiCfg.Singup)
	router.POST("login/", apiCfg.Login)
	router.POST("refresh-token/", apiCfg.RefreshAccessToken)

	// Auth routes
	authRouter.GET("profile/", apiCfg.GetUserProfile)
	authRouter.PATCH("profile/", apiCfg.UpdateUserProfile)
	authRouter.DELETE("profile/", apiCfg.DeleteUserProfile)
	authRouter.PUT("change-password/", apiCfg.ChangePassword)

	return router
}
