package main

import (
	"mikadifo/money-moon/models"
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
	router.GET("/", ping)
	router.GET("/transactions/:bankId", models.GetAllTransactionsByBankId)

	router.Run("localhost:" + PORT)
}
