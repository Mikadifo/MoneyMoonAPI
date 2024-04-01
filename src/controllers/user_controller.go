package controllers

import (
	"context"
	"fmt"
	"mikadifo/money-moon/src/config"
	"mikadifo/money-moon/src/models"
	"mikadifo/money-moon/src/responses"
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
		responses.Send(c, http.StatusBadRequest, responses.ERROR, err.Error())
		return
	}

	if validationError := validate.Struct(&user); validationError != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, validationError.Error())
		return
	}

	newUser, err := getUserByEmail(user.Email)
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	if newUser.Email == user.Email {
		responses.Send(c, http.StatusConflict, responses.ERROR, "User already exists with the following email: "+user.Email)
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
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	responses.Send(c, http.StatusCreated, responses.SUCCESS, result)
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, err.Error())
		return
	}

	if body.Email == "" || body.Password == "" {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, "Email and/or is empty.")
		return
	}

	user, err := getUserByEmail(body.Email)
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	if user.Email != body.Email || user.Password != body.Password {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, "Email or password incorrect.")
		return
	}

	responseData := bson.M{
		"token": "TODO:here goes encrypted jwt token",
	}

	responses.Send(c, http.StatusOK, responses.SUCCESS, responseData)
}

func GetUserByEmail(c *gin.Context) {
	var user models.User
	email := c.Param("email")

	user, err := getUserByEmail(email)
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	if user.Email != email {
		responses.Send(c, http.StatusNotFound, responses.ERROR, "User with email "+email+" not found.")
		return
	}

	responses.Send(c, http.StatusOK, responses.SUCCESS, user)
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
