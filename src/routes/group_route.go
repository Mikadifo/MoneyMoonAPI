package routes

import (
	"mikadifo/money-moon/src/controllers"
	"mikadifo/money-moon/src/middleware"

	"github.com/gin-gonic/gin"
)

func GroupRoute(router *gin.Engine) {
	router.GET("/groups", middleware.RequireAuth, controllers.GetAllGroups)
	router.PUT("/groups/add/:groupId", middleware.RequireAuth, controllers.AddTransactions)
}
