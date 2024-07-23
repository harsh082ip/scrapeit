package authroutes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authcontrollers "github.com/harsh082ip/scrapeit/authService/authControllers"
)

func AuthRoutes(incommingRoutes *gin.Engine) {
	authRouter := incommingRoutes.Group("/auth")
	{
		authRouter.POST("/signup", authcontrollers.SignUp)
		authRouter.POST("/login", authcontrollers.Login)
		authRouter.GET("/", func(c *gin.Context) {
			// Retrieve the cookie named "example_cookie"
			cookie, err := c.Cookie("jwt_key")
			if err != nil {
				// If the cookie does not exist, handle the error
				c.JSON(http.StatusNotFound, gin.H{"message": "Cookie not found"})
				return
			}

			// If the cookie is found, return its value
			c.JSON(http.StatusOK, gin.H{"cookie_value": cookie})
		})
	}
}
