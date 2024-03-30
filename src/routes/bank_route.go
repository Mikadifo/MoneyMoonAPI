package routes

import (
	"mikadifo/money-moon/src/controllers"

	"github.com/gin-gonic/gin"
)

func BankRoute(router *gin.Engine) {
	router.POST("/bank/create", controllers.CreateBank)
	router.GET("/bank/:bankId", controllers.GetBankByID)
}
