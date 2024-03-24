package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Hello!"})
}

func main() {
	router := gin.Default()
	router.GET("/", ping)

	router.Run("localhost:8080")
}
