package controllers

import (
	"context"
	"math"
	"mikadifo/money-moon/src/config"
	"mikadifo/money-moon/src/models"
	"mikadifo/money-moon/src/responses"
	"net/http"
	"strconv"
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

func GetTransactionsByBankId(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	bankId := c.Param("bankId")
	pageQuery := c.DefaultQuery("page", "1")
	limitQuery := c.DefaultQuery("limit", "10")
	var transactions []models.Transaction
	var pages int64
	filter := bson.M{"bankId": bankId}
	defer cancel()

	page, err := strconv.ParseInt(pageQuery, 10, 64)
	if err != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, err.Error())
		return
	}
	limit, err := strconv.ParseInt(limitQuery, 10, 64)
	if err != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, err.Error())
		return
	}

	count, err := transactionsCollection.CountDocuments(ctx, filter)
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}
	pages = int64(math.Ceil(float64(count) / float64(limit)))

	if page < 1 || page > pages || limit < 1 {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, "Page number not found or limit is not a positive number")
		return
	}

	skip := int64((page - 1) * limit)
	projection := bson.M{"_id": 0}
	findOptions := options.FindOptions{Limit: &limit, Skip: &skip}
	findOptions.SetProjection(projection)
	cursor, err := transactionsCollection.Find(ctx, filter, &findOptions)
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

	responseData := bson.M{
		"paginator": bson.M{
			"page":  page,
			"pages": pages,
		},
		"transactions": transactions,
	}

	responses.Send(c, http.StatusOK, responses.SUCCESS, responseData)
}
