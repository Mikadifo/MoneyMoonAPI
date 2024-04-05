package middleware

import (
	"fmt"
	"mikadifo/money-moon/src/controllers"
	"mikadifo/money-moon/src/models"
	"mikadifo/money-moon/src/responses"
	"mikadifo/money-moon/src/utily"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	tokenString := c.GetHeader("access-token")
	if tokenString == "" {
		responses.Unauthorized(c)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secret := utily.GetEnvVar("SECRET")
		return []byte(secret), nil
	})

	if token == nil || err != nil || !token.Valid {
		responses.Unauthorized(c)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			responses.Unauthorized(c)
			return
		}

		var user models.User
		userId := claims["sub"].(string)
		user, err := controllers.GetUserByID(userId)
		if err != nil {
			responses.Unauthorized(c)
			return
		}

		if user.Id.Hex() != userId || user.Email == "" {
			responses.Unauthorized(c)
			return
		}

		c.Next()
	} else {
		responses.Unauthorized(c)
	}
}
