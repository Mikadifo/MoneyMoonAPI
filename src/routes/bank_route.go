package routes

import (
	"mikadifo/money-moon/src/controllers"
	"mikadifo/money-moon/src/middleware"

	"github.com/gin-gonic/gin"
)

func BankRoute(router *gin.Engine) {
	router.POST("/bank/create", middleware.RequireAuth, controllers.CreateBank)
	router.GET("/bank/:bankId", controllers.GetBankByID) //TODO: check if is used or not
}
