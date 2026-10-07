package middleware

import (
	"net/http"
	"strings"

	"github.com/TuanNghia295/BE-FIT/internal/infrastructure/postgres/auth"
	"github.com/gin-gonic/gin"
)

// Create AuthenticatedUserIDKey to store the authenticated user ID in the context
// why need it?
// Because we need to pass the authenticated user ID to the next handler in the chain
// So we can use it to the next handler to get the authenticated user ID from the context
const AuthenticatedUserIDKey = "authenticated_user_id"

func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the Authorization header from the request
		parts := strings.Fields(ctx.GetHeader("Authorization"))
		// Parts will return an array of strings, where the first element is the "Bearer"
		// and the second element is the actual token.
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid bearer token",
			})
			return
		}
		claims, err := jwtService.ValidateAccessToken(parts[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		ctx.Set(AuthenticatedUserIDKey, claims.UserID)
		ctx.Next()
	}

}
