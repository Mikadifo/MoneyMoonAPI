package utily

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVar(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Print("Error oppening .env file")
	}

	return os.Getenv(key)
}
