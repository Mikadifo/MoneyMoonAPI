package responses

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Unauthorized(c *gin.Context) {
	response := DefaultResponse{
		Status:  http.StatusUnauthorized,
		Message: ERROR.String(),
		Data:    "Authorization token is not provided, invalid or expired",
	}
	c.AbortWithStatusJSON(response.Status, response)
}
