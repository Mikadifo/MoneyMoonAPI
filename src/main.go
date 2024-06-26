package main

import (
	"mikadifo/money-moon/src/config"
	"mikadifo/money-moon/src/middleware"
	"mikadifo/money-moon/src/routes"
	"mikadifo/money-moon/src/utily"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello!"})
}

func main() {
	PORT := utily.GetEnvVar("PORT")

	router := gin.Default()
	router.Use(middleware.LocalCors())
	config.ConnectDB()

	router.GET("/", ping)
	routes.UserRoute(router)
	routes.TransactionRoute(router)
	routes.BankRoute(router)
	routes.DebtRoute(router)
	routes.GroupRoute(router)

	router.Run("localhost:" + PORT)
}
