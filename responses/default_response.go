package responses

import (
	"github.com/gin-gonic/gin"
)

type Message int

const (
	SUCCESS Message = iota
	ERROR
)

type DefaultResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (message Message) String() string {
	return [...]string{"success", "error"}[message]
}

func Send(c *gin.Context, status int, message Message, data interface{}) {
	response := DefaultResponse{
		Status:  status,
		Message: message.String(),
		Data:    data,
	}
	c.IndentedJSON(status, response)
}
