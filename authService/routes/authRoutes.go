package routes

import (
	"github.com/gin-gonic/gin"
	authcontrollers "github.com/harsh082ip/scrapeit/authService/authControllers"
)

func AuthRoutes(incommingRoutes *gin.RouterGroup) {
	incommingRoutes.POST("/signup", authcontrollers.SignUp)
}
