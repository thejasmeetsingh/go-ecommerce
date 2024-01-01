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

	router := engine.Group("/api/v1/")
	router.Use(JWTAuth((apiCfg)))
	router.POST("order/", apiCfg.CreateOrder)
	router.GET("order/", apiCfg.GetOrders)
	router.GET("order/:id/", apiCfg.GetOrderDetail)
	router.DELETE("order/:id/", apiCfg.DeleteOrder)
}
