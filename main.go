package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func getEnvVar(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Hello!"})
}

func main() {
	router := gin.Default()
	router.GET("/", ping)

	router.Run("localhost:" + getEnvVar("PORT"))
}
