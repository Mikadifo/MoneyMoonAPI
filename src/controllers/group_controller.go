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
	"go.mongodb.org/mongo-driver/mongo/options"
)

var groupCollection *mongo.Collection = config.GetCollection(config.MongoClient, "Groups")

func CreateGroup(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var body struct {
		Name string `json:"name,omitempty" validate:"required"`
	}
	group := models.Group{}
	userId, exists := c.Get("userId")
	defer cancel()

	if !exists || userId == "" {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, "We couldn't find user's ID in the token.")
		return
	}

	if err := c.BindJSON(&body); err != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, err.Error())
		return
	}

	if validationError := validate.Struct(&body); validationError != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, validationError.Error())
		return
	}

	groupExists, err := groupExists(body.Name)
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	if groupExists {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, "Group already exists with name "+body.Name)
		return
	}

	group.Name = body.Name
	group.UserId = userId.(string)
	group.Total = 0
	group.Transactions = []string{}

	_, err = groupCollection.InsertOne(ctx, group)
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	responses.Send(c, http.StatusCreated, responses.SUCCESS, "Group saved succesfully")
}

func GetAllGroups(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId, exists := c.Get("userId")
	var groups []models.Group
	defer cancel()

	if !exists || userId == "" {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, "We couldn't find user's ID in the token.")
		return
	}

	projection := bson.M{"userId": 0}
	findOptions := options.Find().SetProjection(projection)
	cursor, err := groupCollection.Find(ctx, bson.M{"userId": userId}, findOptions)
	if err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	if err = cursor.All(ctx, &groups); err != nil {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	responses.Send(c, http.StatusCreated, responses.SUCCESS, groups)
}

func AddTransactions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId, exists := c.Get("userId")
	groupIdHex := c.Param("groupId")
	var transactions []string
	var group models.Group
	defer cancel()

	if !exists || userId == "" {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, "We couldn't find user's ID in the token.")
		return
	}

	groupId, err := primitive.ObjectIDFromHex(groupIdHex)
	if err != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, err.Error())
		return
	}

	if err := c.BindJSON(&transactions); err != nil {
		responses.Send(c, http.StatusBadRequest, responses.ERROR, err.Error())
		return
	}

	update := bson.M{"$addToSet": bson.M{"transactions": bson.M{"$each": transactions}}}
	err = groupCollection.FindOneAndUpdate(ctx, bson.M{"_id": groupId}, update).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			responses.Send(c, http.StatusNotFound, responses.ERROR, "Group not found.")
			return

		}

		responses.Send(c, http.StatusInternalServerError, responses.ERROR, err.Error())
		return
	}

	responses.Send(c, http.StatusOK, responses.SUCCESS, "Transactions added succesfully")
}

func groupExists(name string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var group models.Group
	defer cancel()

	err := groupCollection.FindOne(ctx, bson.M{"name": name}).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}

		return true, err
	}

	return true, nil
}
