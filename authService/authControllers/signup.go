package authcontrollers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	models "github.com/harsh082ip/scrapeit/Models"
	"github.com/harsh082ip/scrapeit/db"
	"github.com/harsh082ip/scrapeit/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignUp(c *gin.Context) {

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		// Return a bad request response if there's an error in binding/validation
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Error in request body",
			"error":  err.Error(),
		})
		return
	}

	collName := "Users"
	coll := db.OpenCollection(db.Client, collName)

	emailFiler := bson.M{"email": user.Email}
	usernameFilter := bson.M{"username": user.Username}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	count, err := coll.CountDocuments(ctx, emailFiler)
	if err != nil && err != mongo.ErrNilDocument {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error While Checking for Doc",
			"error":  err.Error(),
			"count":  count,
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Doc Duplication not allowed",
			"error":  "this email already exists",
		})
		return
	}

	// basic email verification check
	_, err = helpers.VerifyEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Email Verification Failed",
			"error":  err.Error(),
		})
		return
	}

	count, err = coll.CountDocuments(ctx, usernameFilter)
	if err != nil && err != mongo.ErrNilDocument {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error While Checking for Doc",
			"error":  err.Error(),
			"count":  count,
		})
		return
	}

	if count > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Doc Duplication not allowed",
			"error":  "this username already exists",
		})
		return
	}

	user.Password, err = helpers.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error in generating hash for password",
			"error":  err.Error(),
		})
		return
	}

	user.ID = primitive.NewObjectID()

	// If all checks are passed we'll store details of the user on a temporary basis in redis

	rdb := db.RedisConnect()
	key := "user:" + user.Email

	// Marshal the user struct
	jsonData, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error in Marshalling body",
			"error":  err.Error(),
		})
		return
	}

	_, err = rdb.Set(ctx, key, jsonData, time.Minute*15).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Unexpected error while Signing up :/",
			"error":  "Cannot Set user data to redis, " + err.Error(),
		})
		return
	}

	// now we need to push the user email to a queue and the worker will pick from there
	key = "emails_to_verify"

	_, err = rdb.LPush(ctx, key, []byte(user.Email)).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Unexpected error while signing up :/",
			"error":  "Cannot push email to queue, " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Ok",
		"msg":    "A verification email has been sent to you. \nIf not received yet, consider checking the spam folder or wait for a couple of minutes :)",
	})

}
