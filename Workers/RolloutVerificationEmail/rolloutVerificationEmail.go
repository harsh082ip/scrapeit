package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/harsh082ip/scrapeit/db"
	"github.com/redis/go-redis/v9"
)

const (
	redisKey           = "emails_to_verify"
	emailServiceURL    = "http://localhost:8002/email/user/"
	sleepDuration      = 2 * time.Second
	errorSleepDuration = 1 * time.Second
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// capture interrupt signal to allow graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		<-sigCh
		log.Println("Received Interrupt signal, shutting down...")
		cancel()
	}()

	RolloutVerificationEmail(ctx)
}

func RolloutVerificationEmail(ctx context.Context) {
	rdb := db.RedisConnect()

	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled, stopping RolloutVerificationEmail")
			return
		default:
			email, err := rdb.RPop(ctx, redisKey).Result()
			log.Println(email)
			if err != nil && err != redis.Nil {
				log.Println("Error in popping emails from the Redis server:", err)
				time.Sleep(errorSleepDuration)
				continue
			}
			if err == redis.Nil {
				log.Println("No emails are in the queue for verification")
				time.Sleep(sleepDuration)
				continue
			}

			success, err := SendEmail(ctx, email)
			if err != nil {
				log.Println("Error in sending email:", err)
				time.Sleep(errorSleepDuration)
				continue
			}

			if !success {
				log.Println("Cannot send an email due to an unexpected error")
			}

			log.Println("Verification email sent to:", email)
		}
	}
}

func SendEmail(ctx context.Context, email string) (bool, error) {
	url := emailServiceURL + email

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("error creating the GET request to send-email service: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("error making the GET request to send-email service: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("error sending an email, code: %v, response body: %v", resp.StatusCode, resp.Body)
	}

	return true, nil
}
