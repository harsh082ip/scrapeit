package authcontrollers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	models "github.com/harsh082ip/scrapeit/Models"
	"github.com/harsh082ip/scrapeit/db"
	authhelper "github.com/harsh082ip/scrapeit/helpers/authHelper"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func Login(c *gin.Context) {

	var user models.LoginUser
	collName := "Users"
	coll := db.OpenCollection(db.Client, collName)
	ctx, cancel := context.WithTimeout(c, time.Second*15)
	defer cancel()

	if err := c.ShouldBindJSON(&user); err != nil {
		// Return a bad request response if there's an error in binding/validation
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "Error in request body",
			"error":  err.Error(),
		})
		return
	}

	// find the user from the db with the given email
	var result models.User
	err := coll.FindOne(ctx, bson.M{"email": user.LoginID}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "No such user found with this email",
				"error":  err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Error occured while searching for the user",
				"error":  err.Error(),
			})
			return
		}
	}

	err = authhelper.ComparePassword(result.Password, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "Password Incorrect!",
			"error":  err.Error(),
		})
		return
	}

	// Generate JWT token for the authenticated user
	userJwtToken, err := authhelper.GenerateJwtToken(result.Email)
	if err != nil {
		// Handle error in generating JWT token
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Error in Generating token",
			"error":  err.Error(),
		})
		return
	}

	// Return successful login response with JWT token
	c.JSON(http.StatusOK, gin.H{
		"status":    "All Good! User Login Successful",
		"Jwt_Token": userJwtToken,
		"user":      result,
	})
}
