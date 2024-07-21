package authhelper

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	models "github.com/harsh082ip/scrapeit/Models"
	"github.com/joho/godotenv"
)

func GenerateJwtToken(loginID string) (string, error) {

	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Println("Error in Loading env :/")
	}

	JWT_SECRET_KEY := os.Getenv("JWT_SECRET_KEY")

	if JWT_SECRET_KEY == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY cannot be empty :/")
	}

	// expiration time = 30 days
	expirationTime := time.Now().Add(time.Hour * 720)

	// creating claims
	claims := &models.Claims{
		CompanyName: "Scrapeit",
		LoginID:     loginID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte(JWT_SECRET_KEY)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
