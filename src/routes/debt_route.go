package routes

import (
	"mikadifo/money-moon/src/controllers"
	"mikadifo/money-moon/src/middleware"

	"github.com/gin-gonic/gin"
)

func DebtRoute(router *gin.Engine) {
	router.GET("/debts", middleware.RequireAuth, controllers.GetUnpaidDebts)
	router.POST("/debts/create", middleware.RequireAuth, controllers.CreateDebt)
	router.PUT("/debts/pay", middleware.RequireAuth, controllers.PayAmount)
	router.DELETE("/debts/remove/:debtName", middleware.RequireAuth, controllers.DeleteDebt)
}
