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

func CreateTransactions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var body []models.Transaction
	var transactions = []interface{}{}
	var invalidTransactions []models.Transaction
	var errors []string
	defer cancel()

	if err := c.BindJSON(&body); err != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, err.Error())
		return
	}

	for _, transaction := range body {
		if validationError := validate.Struct(&transaction); validationError != nil {
			invalidTransactions = append(invalidTransactions, transaction)
			errors = append(errors, validationError.Error())
			continue
		}

		transactions = append(transactions, transaction)
	}

	responseData := bson.M{
		"failed": invalidTransactions,
		"errors": errors,
	}
	if len(transactions) == 0 {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, responseData)
		return
	}

	result, err := transactionsCollection.InsertMany(ctx, transactions)
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	responseData = bson.M{
		"insertedCount": len(result.InsertedIDs),
		"failed":        invalidTransactions,
		"failedCount":   len(invalidTransactions),
		"errors":        errors,
	}
	if len(invalidTransactions) > 0 {
		responses.Send(c, http.StatusMultiStatus, responses.PARTIAL, responseData)
		return
	}

	responses.Send(c, http.StatusCreated, responses.SUCCESS, bson.M{"insertedCount": len(result.InsertedIDs)})
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
