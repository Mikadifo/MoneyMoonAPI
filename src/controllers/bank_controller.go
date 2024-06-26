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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var bankCollection *mongo.Collection = config.GetCollection(config.MongoClient, "Banks")

func CreateBank(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var body models.Bank
	var banks []models.Bank
	defer cancel()

	if err := c.BindJSON(&body); err != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, err.Error())
		return
	}

	if validationError := validate.Struct(&body); validationError != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, validationError.Error())
		return
	}

	cursor, err := bankCollection.Find(ctx, bson.M{"userId": body.UserId})
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	if err = cursor.All(ctx, &banks); err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	for _, bankObj := range banks {
		if bankObj.Name == body.Name {
			responses.Send(c, http.StatusConflict, responses.ERROR, "Bank with name "+body.Name+" already exists")
			return
		}
	}

	result, err := bankCollection.InsertOne(ctx, body)
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	user, err := GetUserByID(body.UserId)
	if err != nil || user.Email == "" {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	bankId := result.InsertedID
	stringBankId := bankId.(primitive.ObjectID).Hex()
	user.Banks = append(user.Banks, stringBankId)
	filter := bson.M{"_id": user.Id}
	update := bson.M{"$set": bson.M{"banks": user.Banks}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	responses.Send(c, http.StatusCreated, responses.SUCCESS, result)
}

func GetBankByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var bank models.Bank
	bankId := c.Param("bankId")
	defer cancel()

	objId, err := primitive.ObjectIDFromHex(bankId)
	if err != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, "Bank id is not valid.")
		return
	}

	err = bankCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&bank)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			responses.Send(c, http.StatusNotFound, responses.ERROR, "Bank not found.")
			return

		}
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	responses.Send(c, http.StatusOK, responses.SUCCESS, bank)
}
