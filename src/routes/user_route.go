package routes

import (
	"mikadifo/money-moon/src/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.POST("/signup", controllers.CreateUser)
	router.POST("/login", controllers.Login)
	router.GET("/user/:email", controllers.GetUserByEmail)
}
