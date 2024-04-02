package routes

import (
	"mikadifo/money-moon/src/controllers"

	"github.com/gin-gonic/gin"
)

func TransactionRoute(router *gin.Engine) {
	router.GET("/transactions/:bankId", controllers.GetAllTransactionsByBankId)
	router.POST("/transactions/create", controllers.CreateTransaction)
}
