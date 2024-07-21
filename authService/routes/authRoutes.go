package authroutes

import (
	"github.com/gin-gonic/gin"
	authcontrollers "github.com/harsh082ip/scrapeit/authService/authControllers"
)

func AuthRoutes(incommingRoutes *gin.Engine) {
	authRouter := incommingRoutes.Group("/auth")
	{
		authRouter.POST("/signup", authcontrollers.SignUp)
		authRouter.POST("/login", authcontrollers.Login)
	}
}
