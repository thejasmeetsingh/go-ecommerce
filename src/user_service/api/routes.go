package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRoutes(engine *gin.Engine, apiCfg *APIConfig) {
	// Default route for health check
	engine.GET("/health-check/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Up and Running!",
		})
	})

	// Public API Routes
	pubRouter := engine.Group("/api/v1/")
	authPubRouter := pubRouter.Group("")
	authPubRouter.Use(JWTAuth((apiCfg)))

	// Non auth routes
	pubRouter.POST("register/", apiCfg.Singup)
	pubRouter.POST("login/", apiCfg.Login)
	pubRouter.POST("refresh-token/", apiCfg.RefreshAccessToken)

	// Auth routes
	authPubRouter.GET("profile/", apiCfg.GetUserProfile)
	authPubRouter.PATCH("profile/", apiCfg.UpdateUserProfile)
	authPubRouter.DELETE("profile/", apiCfg.DeleteUserProfile)
	authPubRouter.PUT("change-password/", apiCfg.ChangePassword)

	// Internal API Routes
	pvtRouter := engine.Group("/internal/v1/")
	pvtRouter.Use(InternalAPIAuth(apiCfg))

	pvtRouter.POST("token/", apiCfg.GetUserFromToken)
}
