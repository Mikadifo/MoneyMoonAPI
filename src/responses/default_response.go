package responses

import (
	"github.com/gin-gonic/gin"
)

type DefaultResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Send(c *gin.Context, status int, message Message, data interface{}) {
	response := DefaultResponse{
		Status:  status,
		Message: message.String(),
		Data:    data,
	}
	c.JSON(status, response)
}
