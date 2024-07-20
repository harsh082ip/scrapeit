package emailroutes

import (
	"github.com/gin-gonic/gin"
	emailcontrollers "github.com/harsh082ip/scrapeit/emailService/emailControllers"
)

func EmailRoutes(incommingRoutes *gin.Engine) {
	routes := incommingRoutes.Group("/email")
	{
		routes.GET("/user/:email", emailcontrollers.SendEmailToVerify)
		routes.GET("/verify/:email", emailcontrollers.VerifyUserEmail)
	}
}
