package routes

import (
	"mikadifo/money-moon/src/controllers"
	"mikadifo/money-moon/src/middleware"

	"github.com/gin-gonic/gin"
)

func GroupRoute(router *gin.Engine) {
	router.GET("/groups", middleware.RequireAuth, controllers.GetAllGroups)
	router.GET("/groups/:groupId", middleware.RequireAuth, controllers.GetTransactionsByGroupId)
	router.POST("/groups", middleware.RequireAuth, controllers.CreateGroup)
	router.PUT("/groups/add/:groupId", middleware.RequireAuth, controllers.AddTransactions)
	router.DELETE("/groups/delete/:groupId", middleware.RequireAuth, controllers.DeleteGroupTransaction)
}
