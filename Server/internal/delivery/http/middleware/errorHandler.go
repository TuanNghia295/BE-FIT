package middleware

//  General Error Handling
import (
	"errors"
	"net/http"

	domainErrors "github.com/TuanNghia295/BE-FIT/internal/domains/errors"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.Next()

		if len(ctx.Errors) == 0 {
			return
		}

		err := ctx.Errors.Last().Err

		switch {
		case errors.Is(err, domainErrors.ErrUserNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"code":    "USER_NOT_FOUND",
				"message": "User not found",
			})

		case errors.Is(err, domainErrors.ErrUserAlreadyExists):
			ctx.JSON(http.StatusConflict, gin.H{
				"success": false,
				"code":    "USER_ALREADY_EXISTS",
				"message": "User already exists",
			})

		case errors.Is(err, domainErrors.ErrUnauthorized):
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"code":    "UNAUTHORIZED",
				"message": "Unauthorized",
			})

		default:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Internal server error",
			})
		}
	}
}
