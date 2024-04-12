package routes

import (
	"mikadifo/money-moon/src/controllers"
	"mikadifo/money-moon/src/middleware"

	"github.com/gin-gonic/gin"
)

func DebtRoute(router *gin.Engine) {
	router.POST("/debts/create", middleware.RequireAuth, controllers.CreateDebt)
	router.PUT("/debts/pay", middleware.RequireAuth, controllers.PayAmount)
}
