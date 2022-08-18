package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mustafabalila/golang-api/config"
)

type NotifyInput struct {
	Token string
	Title string
	Body  string
}

func NotifyUser(input NotifyInput) error {
	cfg := config.GetConfig()

	fcmNotification := map[string]interface{}{
		"to": input.Token,
		"notification": map[string]string{
			"title": input.Title,
			"body":  input.Body,
		},
	}

	bodyStr, err := json.Marshal(fcmNotification)
	if err != nil {
		return err
	}

	bodyByes := []byte(bodyStr)
	bodyReader := bytes.NewReader(bodyByes)

	req, err := http.NewRequest(http.MethodPost, "https://fcm.googleapis.com/fcm/send", bodyReader)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("key=%s", cfg.FCMAPIKey))

	client := http.Client{Timeout: 60 * time.Second}
	_, err = client.Do(req)

	if err != nil {
		return err
	}

	return nil
}

func NotifyUsers(tokens []string, notification NotifyInput) error {
	for _, token := range tokens {
		err := NotifyUser(NotifyInput{Token: token, Title: notification.Title, Body: notification.Body})
		if err != nil {
			return err
		}
	}
	return nil
}
