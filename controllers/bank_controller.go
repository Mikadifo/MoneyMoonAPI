package controllers

import (
	"context"
	"mikadifo/money-moon/config"
	"mikadifo/money-moon/models"
	"mikadifo/money-moon/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var bankCollection *mongo.Collection = config.GetCollection(config.MongoClient, "Banks")

func CreateBank(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var bank models.Bank
	var banks []models.Bank
	defer cancel()

	if err := c.BindJSON(&bank); err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.DefaultResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	if validationError := validate.Struct(&bank); validationError != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.DefaultResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationError.Error()}})
		return
	}

	cursor, err := bankCollection.Find(ctx, bson.M{"userId": bank.UserId})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.DefaultResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	if err = cursor.All(ctx, &banks); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.DefaultResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	for _, bankObj := range banks {
		if bankObj.Name == bankObj.Name {
			c.IndentedJSON(http.StatusConflict, responses.DefaultResponse{Status: http.StatusConflict, Message: "error", Data: map[string]interface{}{"data": "Bank with name " + bank.Name + " already exists"}})
			return
		}
	}

	result, err := bankCollection.InsertOne(ctx, bank)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.DefaultResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	c.IndentedJSON(http.StatusCreated, responses.DefaultResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
}

func GetBankByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var bank models.Bank
	bankId := c.Param("bankId")
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(bankId)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, responses.DefaultResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Bank id is not valid."}})
		return
	}

	err = bankCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&bank)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.IndentedJSON(http.StatusNotFound, responses.DefaultResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "Bank not found."}})
			return

		}
		c.IndentedJSON(http.StatusInternalServerError, responses.DefaultResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	c.IndentedJSON(http.StatusOK, responses.DefaultResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": bank}})
}
