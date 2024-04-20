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

var groupCollection *mongo.Collection = config.GetCollection(config.MongoClient, "Groups")

func GetAllGroups(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	userId, exists := c.Get("userId")
	var groups []models.Group
	defer cancel()

	if !exists || userId == "" {
		responses.Send(c, http.StatusInternalServerError, responses.ERROR, "We couldn't find user's ID in the token.")
		return
	}

	projection := bson.M{"_id": 0, "userId": 0}
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
