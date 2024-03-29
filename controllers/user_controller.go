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
	"go.mongodb.org/mongo-driver/bson"
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

	newUser, err := getUserByEmail(user.Email)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	if newUser.Email == user.Email {
		c.IndentedJSON(http.StatusConflict, responses.UserResponse{Status: http.StatusConflict, Message: "error", Data: map[string]interface{}{"data": "User already exists with the following email: " + user.Email}})
		return
	}

	newUser = models.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Banks:    []string{},
		Debts:    []models.Debt{},
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	c.IndentedJSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
}

func GetUserByEmail(c *gin.Context) {
	var user models.User
	email := c.Param("email")

	user, err := getUserByEmail(email)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	if user.Email != email {
		c.IndentedJSON(http.StatusNotFound, responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with email " + email + " not found."}})
		return
	}

	c.IndentedJSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": user}})
}

func getUserByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, nil
		}

		return models.User{}, err
	}

	return user, nil
}