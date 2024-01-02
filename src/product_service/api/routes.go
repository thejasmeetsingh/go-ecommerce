package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func GetRoutes(engine *gin.Engine, apiCfg *APIConfig) {
	// Default route for health check
	engine.GET("/health-check/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Up and Running!",
		})
	})

	// Add Rate limiter middleware
	engine.Use(RateLimiter(apiCfg))

	// Add prometheus middleware and route
	engine.Use(PrometheusMiddleware())
	engine.GET("/metrics/", gin.WrapH(promhttp.Handler()))

	pubRouter := engine.Group("/api/v1/")
	pubRouter.Use(JWTAuth(apiCfg))
	pubRouter.POST("product/", apiCfg.CreateProduct)
	pubRouter.GET("product/:id/", apiCfg.GetProductDetails)
	pubRouter.PATCH("product/:id/", apiCfg.UpdateProductDetails)
	pubRouter.DELETE("product/:id/", apiCfg.DeleteProduct)
	pubRouter.GET("product/", apiCfg.GetProducts)

	pvtRouter := engine.Group("/internal/v1/")
	pvtRouter.Use(InternalAPIAuth(apiCfg))
	pvtRouter.POST("product-details/", apiCfg.GetProductIDToDetails)
}
