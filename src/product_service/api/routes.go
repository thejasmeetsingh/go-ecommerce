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
	router.Use(Auth(apiCfg))
	router.POST("product/", apiCfg.CreateProduct)
	router.GET("product/:id/", apiCfg.GetProductDetails)
	router.PATCH("product/:id/", apiCfg.UpdateProductDetails)
	router.DELETE("product/:id/", apiCfg.DeleteProduct)
	router.GET("product/", apiCfg.GetProducts)
}
