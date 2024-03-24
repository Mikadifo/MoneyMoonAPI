package models

import (
	"context"
	"mikadifo/money-moon/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type Transaction struct {
	_id         string  `json:"_id"`
	bankId      string  `json:"bankId"`
	transType   string  `json:"type"`
	description string  `json:"description"`
	amount      float32 `json:"amount"`
	date        string  `json:"date"`
}

func GetAllTransactionsByBankId(c *gin.Context) {
	bankId := c.Param("bankId")

	cursor, err := config.MongoClient.Database("MoneyMoon").Collection("Transactions").Find(context.TODO(), bson.D{{Key: "bankId", Value: bankId}})

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var transactions []bson.M
	if err = cursor.All(context.TODO(), &transactions); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, transactions)
}
