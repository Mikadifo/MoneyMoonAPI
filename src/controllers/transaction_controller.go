package controllers

import (
	"context"
	"mikadifo/money-moon/src/config"
	"mikadifo/money-moon/src/models"
	"mikadifo/money-moon/src/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var transactionsCollection *mongo.Collection = config.GetCollection(config.MongoClient, "Transactions")

func CreateTransaction(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var transaction models.Transaction
	defer cancel()

	if err := c.BindJSON(&transaction); err != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, err.Error())
		return
	}

	if validationError := validate.Struct(&transaction); validationError != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, validationError.Error())
		return
	}

	result, err := transactionsCollection.InsertOne(ctx, transaction)
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	responses.Send(c, http.StatusCreated, responses.SUCCESS, result)
}

func GetAllTransactionsByBankId(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	bankId := c.Param("bankId")
	var transactions []models.Transaction
	defer cancel()

	projection := bson.M{"_id": 0}
	opts := options.Find().SetProjection(projection)
	cursor, err := transactionsCollection.Find(ctx, bson.M{"bankId": bankId}, opts)
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	if err = cursor.All(ctx, &transactions); err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	if transactions == nil {
		transactions = []models.Transaction{}
	}

	responses.Send(c, http.StatusOK, responses.SUCCESS, transactions)
}
