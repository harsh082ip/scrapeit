package emailcontrollers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	models "github.com/harsh082ip/scrapeit/Models"
	"github.com/harsh082ip/scrapeit/db"
	"github.com/harsh082ip/scrapeit/helpers"
	"github.com/redis/go-redis/v9"
	"gopkg.in/mgo.v2/bson"
)

func VerifyUserEmail(c *gin.Context) {

	email := c.Param("email")
	rdb := db.RedisConnect()
	key := "user:" + email
	ctx, cancel := context.WithTimeout(c, time.Second*15)
	defer cancel()
	var user models.User
	collName := "Users"
	coll := db.OpenCollection(db.Client, collName)
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

	// finally insert user data to db
	_, err = coll.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error is attempting to SignUp",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "User SignUp Successful",
	})
}
