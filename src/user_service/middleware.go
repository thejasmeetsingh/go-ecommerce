package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thejasmeetsingh/go-ecommerce/src/user_service/handlers"
	"github.com/thejasmeetsingh/go-ecommerce/src/user_service/utils"
)

// Validate the request by checking wheather or not they have the valid JWT access token or not
//
// Token format: Bearer <TOKEN>
func JWTAuth(apiCfg *handlers.ApiConfig) gin.HandlerFunc {
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

		// Verify the token and get the encoded payload which is the userID string
		claims, err := utils.VerifyToken(authToken[1])
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid authentication token"})
			ctx.Abort()
			return
		}

		// Check the validity of the token
		if !time.Unix(claims.ExpiresAt.Unix(), 0).After(time.Now()) {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Authentication token is expired"})
			ctx.Abort()
			return
		}

		// Convert the userID string to UUID
		userID, err := uuid.Parse(claims.Data)
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid authentication token"})
			ctx.Abort()
			return
		}

		// Fetch the user by the ID
		dbUser, err := apiCfg.Queries.GetUserById(ctx, userID)
		if err != nil {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Something went wrong"})
			ctx.Abort()
			return
		}

		ctx.Set("user", dbUser)

		// Further call the given handler and send the user instance as well
		ctx.Next()
	}
}
