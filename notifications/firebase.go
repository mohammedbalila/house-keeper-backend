package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/mustafabalila/golang-api/config"
	"github.com/mustafabalila/golang-api/db"
)

type NotifyInput struct {
	// Token is the firebase token of the user
	// Title is the title of the notification
	// Body is the body of the notification
	Token string
	Title string
	Body  string
}

// NotifyUserByToken sends a notification to a user by their firebase token
func NotifyUserByToken(input NotifyInput) error {
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

	fcmRoot := "https://fcm.googleapis.com/fcm/send"
	req, err := http.NewRequest(http.MethodPost, fcmRoot, bodyReader)

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

// NotifyUsersByToken sends a notification to a list of users by their firebase token
func NotifyUsersByToken(tokens []string, notification NotifyInput) error {
	for _, token := range tokens {
		err := NotifyUserByToken(NotifyInput{Token: token, Title: notification.Title, Body: notification.Body})
		if err != nil {
			return err
		}
	}
	return nil
}

// NotifyUsersWithIds sends a notification to a list of users by their database ids
func NotifyUsersWithIds(ids []string, message string, title string) error {
	var err error
	users := &[]db.User{}
	err = db.Database.Model(users).Where("id IN (?)", pg.In(ids)).Column("firebase_token").Select()
	if err != nil {
		return err
	}
	tokens := []string{}
	for _, user := range *users {
		tokens = append(tokens, user.FirebaseToken)
	}

	err = NotifyUsersByToken(
		tokens,
		NotifyInput{Title: title, Body: message})

	if err != nil {
		return err
	}

	return nil
}
