package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/scrapeit/authService/routes"
)

const (
	WEBPORT = ":8001"
)

func main() {

	router := gin.Default()

	authRouter := router.Group("/auth")
	{
		routes.AuthRoutes(authRouter)
	}

	if err := router.Run(WEBPORT); err != nil {
		log.Fatalf("Error Starting the server on %v, %v", WEBPORT, err)
	}
}
