package routes

import (
	"mikadifo/money-moon/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.POST("/signup", controllers.CreateUser)
	router.GET("/user/:email", controllers.GetUserByEmail)
}
