package main

import (
	"log"

	"github.com/gin-gonic/gin"
	emailroutes "github.com/harsh082ip/scrapeit/emailService/routes"
)

const (
	WEBPORT = ":8002"
)

func main() {

	router := gin.Default()

	emailroutes.EmailRoutes(router)

	if err := router.Run(WEBPORT); err != nil {
		log.Printf("Error in starting the email server on port %v, \nerror: %v", WEBPORT, err)
	}
}
