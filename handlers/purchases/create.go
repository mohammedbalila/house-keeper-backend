package purchases

import (
	"fmt"
	"math"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/db"
	"github.com/mustafabalila/golang-api/notifications"
	"github.com/mustafabalila/golang-api/utils/common"
	"github.com/mustafabalila/golang-api/utils/logger"
)

type createPurchaseInput struct {
	Name            string   `json:"name"`
	TotalPrice      float64  `json:"totalPrice"`
	Description     string   `json:"description"`
	Category        int      `json:"category"`
	PaymentProgress int      `json:"paymentProgress"`
	Subscribers     []string `json:"subscribers"`
}

// CreatePurchase creates a new purchase it also creates
// a new purchase subscription fot the users in the subscribers array
// and notifies them it accepts a json body mapped as createPurchaseInput
func CreatePurchase(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var _, err error
	input := &createPurchaseInput{}
	err = c.Bind(input)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	sharePrice := input.TotalPrice / float64(len(input.Subscribers)+1)

	purchase := &db.Purchase{
		UserId:      fmt.Sprintf("%s", c.Get("userId")),
		Name:        input.Name,
		Category:    input.Category,
		SharePrice:  math.Round(sharePrice),
		Description: input.Description,
		TotalPrice:  math.Round(input.TotalPrice),
	}
	_, err = db.Database.Model(purchase).Insert()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	tx, err := db.Database.Begin()
	if err != nil {
		tx.Rollback()
		logger.Error(err.Error())

		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer tx.Close()

	for _, userId := range input.Subscribers {
		payment := &db.PurchaseSubscription{
			PurchaseId: purchase.Id,
			Status:     common.Statuses[common.CREATED],
			UserId:     userId,
		}
		_, err = tx.Model(payment).Insert()
		if err != nil {
			tx.Rollback()
			logger.Error(e.Error())
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	tx.Commit()
	message := fmt.Sprintf("A new purchase (%s) was made and your share is %.1f", purchase.Name, purchase.SharePrice)
	notifications.NotifyUsersWithIds(input.Subscribers, message, "New Purchase")
	response := map[string]interface{}{
		"purchase": purchase,
	}
	return c.JSON(http.StatusCreated, response)
}
