package controllers

import (
	"context"
	"mikadifo/money-moon/src/models"
	"mikadifo/money-moon/src/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

	userObjId, _ := primitive.ObjectIDFromHex(userId.(string))
	update := bson.M{"$push": bson.M{"debts": debt}}
	result, err := userCollection.UpdateByID(ctx, userObjId, update)

	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	responses.Send(c, http.StatusCreated, responses.SUCCESS, result)
}
