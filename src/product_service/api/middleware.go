package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/thejasmeetsingh/go-ecommerce/product_service/shared"
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

// Middleware for checking if a valid API secret is passed or not
func InternalAPIAuth(apiCfg *APIConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiSecret := ctx.GetHeader("Secret")
		internalAPISecret := os.Getenv("INTERNAL_API_SECRET")

		if internalAPISecret == "" {
			panic("Internal API secret is not configured")
		}

		if apiSecret == "" || apiSecret != internalAPISecret {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid API Secret"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
