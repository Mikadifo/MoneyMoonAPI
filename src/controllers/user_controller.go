package controllers

import (
	"context"
	"mikadifo/money-moon/src/config"
	"mikadifo/money-moon/src/encryption"
	"mikadifo/money-moon/src/models"
	"mikadifo/money-moon/src/responses"
	"mikadifo/money-moon/src/utily"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		Password: encryption.HashPassword(user.Password),
		Banks:    []string{},
		Debts:    []models.Debt{},
	}

	_, err = userCollection.InsertOne(ctx, newUser)
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	responses.Send(c, http.StatusCreated, responses.SUCCESS, bson.M{})
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

	passwordIsValid := encryption.VerifyPassword(user.Password, body.Password)
	if user.Email != body.Email || !passwordIsValid {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, "Email or password incorrect.")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id.Hex(),
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	secret := utily.GetEnvVar("SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		tokenString = ""
	}

	responseData := bson.M{
		"token": tokenString,
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

func GetUserBanks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var banks []models.Bank
	userId, exists := c.Get("userId")
	defer cancel()
	if !exists || userId == "" {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, "We couldn't find user's ID in the token.")
		return
	}

	projection := bson.M{"_id": 1, "name": 1}
	opts := options.Find().SetProjection(projection)
	cursor, err := bankCollection.Find(ctx, bson.M{"userId": userId}, opts)
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	if err = cursor.All(ctx, &banks); err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	responses.Send(c, http.StatusOK, responses.SUCCESS, banks)
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

func GetUserByID(id string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()
	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}

	err = userCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.User{}, nil
		}

		return models.User{}, err
	}

	return user, nil
}
