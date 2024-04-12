package controllers

import (
	"context"
	"mikadifo/money-moon/src/models"
	"mikadifo/money-moon/src/responses"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUnpaidDebts(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists || userId == "" {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, "We couldn't find user's ID in the token.")
		return
	}

	user, err := GetUserByID(userId.(string))
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	if user.Id.Hex() != userId {
		responses.Send(c, http.StatusNotFound, responses.ERROR, "User not found")
		return
	}

	responses.Send(c, http.StatusOK, responses.SUCCESS, user.Debts)
}

func CreateDebt(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var debt models.Debt
	userId, exists := c.Get("userId")
	defer cancel()

	if !exists || userId == "" {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, "We couldn't find user's ID in the token.")
		return
	}

	if err := c.BindJSON(&debt); err != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, err.Error())
		return
	}

	if validationError := validate.Struct(&debt); validationError != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, validationError.Error())
		return
	}

	user, err := GetUserByID(userId.(string))
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	if user.Id.Hex() != userId {
		responses.Send(c, http.StatusNotFound, responses.ERROR, "User not found")
		return
	}

	for _, debtObj := range user.Debts {
		if debtObj.Name == debt.Name {
			responses.Send(c, http.StatusConflict, responses.ERROR, "You already have a debt with this name")
			return
		}
	}

	update := bson.M{"$push": bson.M{"debts": debt}}
	_, err = userCollection.UpdateByID(ctx, user.Id, update)

	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	responses.Send(c, http.StatusCreated, responses.SUCCESS, "Debt successfully created")
}

func PayAmount(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var debt models.Debt
	debtName := c.Query("name")
	amountString := c.Query("amount")
	userId, exists := c.Get("userId")
	defer cancel()

	if debtName == "" || amountString == "" {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, "Debt name or amount not provided")
		return
	}

	amount, err := strconv.ParseFloat(amountString, 64)
	if err != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, "Amount should be a number")
		return
	}

	if !exists || userId == "" {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, "We couldn't find user's ID in the token.")
		return
	}

	user, err := GetUserByID(userId.(string))
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	if user.Id.Hex() != userId {
		responses.Send(c, http.StatusNotFound, responses.ERROR, "User not found")
		return
	}
	userObjId, _ := primitive.ObjectIDFromHex(userId.(string))

	for _, debtObj := range user.Debts {
		if debtObj.Name == debtName {
			debt = debtObj
			break
		}
	}

	if debt.Name == "" {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, "Debt with name '"+debtName+"' not found")
		return
	}

	if debt.Amount-(*debt.Payed+amount) < 0 {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, "Payment amount exceeds remaining debt")
		return
	}

	filter := bson.M{"_id": userObjId, "debts.name": debtName}
	update := bson.M{"$inc": bson.M{"debts.$.payed": amount}}
	_, err = userCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	responses.Send(c, http.StatusOK, responses.SUCCESS, "Debt updated")
}
