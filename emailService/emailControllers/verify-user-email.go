package emailcontrollers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	models "github.com/harsh082ip/scrapeit/Models"
	"github.com/harsh082ip/scrapeit/db"
	"github.com/harsh082ip/scrapeit/helpers"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func VerifyUserEmail(c *gin.Context) {

	email := c.Param("email")
	rdb := db.RedisConnect()
	key := "user:" + email
	ctx, cancel := context.WithTimeout(c, time.Second*30)
	defer cancel()
	var user models.User
	collName := "Users"
	collName2 := "Credits"
	mongoClient := db.Client
	coll := db.OpenCollection(mongoClient, collName)
	coll2 := db.OpenCollection(mongoClient, collName2)
	emailFilter := bson.M{"email": user.Email}

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error in param",
			"error":  "email cannot be empty",
		})
		return
	}

	// we'll still double check basic verification
	_, err := helpers.VerifyEmail(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Basic Verification check failed",
			"error":  err.Error(),
		})
		return
	}

	count, err := coll.CountDocuments(ctx, emailFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error verifying user in database",
			"error":  err.Error(),
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "User already exists",
			"error":  "A user with this email already exists. Please consider logging in.",
		})
		return
	}

	res, err := rdb.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error in getting user data from redis db :/",
			"error":  err.Error(),
		})
		return
	}

	if err == redis.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "INVALID REQUEST!",
			"error":  "There is no existing signup request. Please consider creating one.",
		})
	}

	if err := json.Unmarshal([]byte(res), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error in Unmarshalling the user data",
			"error":  err.Error(),
		})
		return
	}

	defaultCredits := &models.AppCredits{
		Email:        user.Email,
		TotalCredits: 10,
	}

	session, err := mongoClient.StartSession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Failed to start a mongoDB session",
			"error":  err.Error(),
		})
	}
	defer session.EndSession(ctx)

	// Define transaction function
	transactionCallback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		if _, err := coll.InsertOne(ctx, user); err != nil {
			return nil, fmt.Errorf("failed in Creating the User, %v", err)
		}

		if _, err := coll2.InsertOne(ctx, defaultCredits); err != nil {
			return nil, fmt.Errorf("error in setting Credit Score for the user, %v", err)
		}

		return nil, nil
	}

	// Execute the transaction
	result, err := session.WithTransaction(context.TODO(), transactionCallback, options.Transaction())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "DB transaction failed :/",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "User SignUp Successful",
		"result": result,
	})
}

// 36:19
