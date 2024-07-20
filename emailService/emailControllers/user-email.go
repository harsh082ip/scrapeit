package emailcontrollers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/harsh082ip/scrapeit/helpers"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

type EmailData struct {
	ID string `json:"email_id"`
}

func SendEmailToVerify(c *gin.Context) {

	email := c.Param("email")

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

	success, err := SendVerificationEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Cannot Send a verification email",
			"error":  err.Error(),
		})
		return
	}

	if !success {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Sending email verification failed due to an unexpected error",
			"error":  "Unknown :/",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Verification email send successfully. \nIf not received consider checking the spam folder",
	})
}

func SendVerificationEmail(email string) (bool, error) {

	// load env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println(".env failed to load")
	}

	sender := os.Getenv("SENDEREMAIL")
	app_password := os.Getenv("APPPASSWORD")
	fmt.Println(sender, app_password)
	if sender == "" || app_password == "" {
		return false, fmt.Errorf("sender-email and app-password not set via env,\nConsider setting that")
	}

	// Get the Html files
	var body bytes.Buffer
	fmt.Println(os.Getwd())
	temp, err := template.ParseFiles("../../static/email/verify.html")
	if err != nil {
		return false, fmt.Errorf("error in parsing html tempelate, %v", err)
	}
	data := &EmailData{
		ID: email,
	}
	err = temp.Execute(&body, data)
	if err != nil {
		return false, fmt.Errorf("error in modifying the html tempelate, %v", err)
	}

	// send with gomail
	m := gomail.NewMessage()
	m.SetHeader("From", sender)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Scrapeit email Verification")
	m.SetBody("text/html", body.String())
	// m.Attach()

	dialer := gomail.NewDialer("smtp.gmail.com", 587, sender, app_password)

	// send email
	if err := dialer.DialAndSend(m); err != nil {
		return false, fmt.Errorf("error in sending verification email, %v", err)
	}

	return true, nil
}
