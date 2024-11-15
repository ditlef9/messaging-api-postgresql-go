// middlewares/auth.go

package middlewares

import (
	"ekeberg.com/messaging-api-postgresql-go/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	// Check if the Authorization header is present
	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}

	// Remove the "Bearer " prefix from the token if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}
	// Verify the token and extract claims
	userId, humanOrService, err := utils.VerifyToken(token) // Get human_or_service from VerifyToken

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}

	// Set the user ID and human_or_service information in the context
	context.Set("userId", userId)
	context.Set("human_or_service", humanOrService) // Set human_or_service in context
	context.Next()
}
