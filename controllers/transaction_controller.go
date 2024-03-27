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

var transactionsCollection *mongo.Collection = config.GetCollection(config.MongoClient, "Transactions")
var validate = validator.New()

func GetAllTransactionsByBankId(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	bankId := c.Param("bankId")
	var transactions []models.Transaction
	defer cancel()

	cursor, err := transactionsCollection.Find(ctx, bson.M{"bankId": bankId})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.TransactionResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	if err = cursor.All(ctx, &transactions); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, responses.TransactionResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		return
	}

	c.IndentedJSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": transactions}})
}
