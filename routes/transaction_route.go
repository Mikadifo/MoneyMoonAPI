package routes

import (
	"mikadifo/money-moon/controllers"

	"github.com/gin-gonic/gin"
)

func TransactionRoute(router *gin.Engine) {
	router.GET("/transactions/:bankId", controllers.GetAllTransactionsByBankId)
}
