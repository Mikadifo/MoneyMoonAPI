package controllers

import (
	"context"
	"mikadifo/money-moon/config"
	"mikadifo/money-moon/models"
	"mikadifo/money-moon/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = config.GetCollection(config.MongoClient, "Users")
var validate = validator.New()

func CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	if err := c.BindJSON(&user); err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	if validationError := validate.Struct(&user); validationError != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationError.Error()}})
		return
	}

	newUser := models.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Debts:    []models.Debt{},
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	c.IndentedJSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
}
