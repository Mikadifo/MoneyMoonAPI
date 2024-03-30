package main

import (
	"mikadifo/money-moon/config"
	"mikadifo/money-moon/routes"
	"mikadifo/money-moon/utily"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Hello!"})
}

func main() {
	PORT := utily.GetEnvVar("PORT")

	router := gin.Default()
	config.ConnectDB()

	router.GET("/", ping)
	routes.UserRoute(router)
	routes.TransactionRoute(router)
	routes.BankRoute(router)

	router.Run("localhost:" + PORT)
}
