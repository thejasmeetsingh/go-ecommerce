package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/thejasmeetsingh/go-ecommerce/order_service/shared"
)

func JWTAuth(apiCfg *APIConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		headerAuthToken := ctx.GetHeader("Authorization")

		if headerAuthToken == "" {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Authentication required"})
			ctx.Abort()
			return
		}

		// Split the token string
		authToken := strings.Split(headerAuthToken, " ")

		// Validate the token string
		if len(authToken) != 2 || authToken[0] != "Bearer" {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid authentication format"})
			ctx.Abort()
			return
		}

		// Fetch userID and set to the context
		userID, err := shared.GetUserFromToken(apiCfg.Cache, ctx, authToken[1])
		if err != nil {
			log.Errorln(err)
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Something went"})
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}

// Prometheus middleware to record HTTP request timings
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Continue processing the request
		c.Next()

		// Collect metrics after the request is processed
		duration := time.Since(startTime).Seconds()

		httpRequestsTotal := GetPromRequestTotal()
		httpRequestDuration := GetPromRequestDuration()

		httpRequestsTotal.WithLabelValues(c.FullPath(), c.Request.Method).Inc()
		httpRequestDuration.WithLabelValues(c.FullPath(), c.Request.Method).Observe(duration)
	}
}
