package routes

import (
	"mikadifo/money-moon/src/controllers"
	"mikadifo/money-moon/src/middleware"

	"github.com/gin-gonic/gin"
)

func TransactionRoute(router *gin.Engine) {
	router.POST("/transactions/create", controllers.CreateTransaction)
	router.GET("/transactions/:bankId", middleware.RequireAuth, controllers.GetAllTransactionsByBankId)
}
