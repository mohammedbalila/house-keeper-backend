package users

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/db"
	"github.com/mustafabalila/golang-api/utils/common"
	"github.com/mustafabalila/golang-api/utils/logger"
)

// GetUserPayments list the payments (purchase subscriptions) of a user.
// It accepts a category and a date range as query parameters
func GetUserPayments(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var err error
	var count int
	category := c.QueryParams().Get("category")
	dateStr := c.QueryParams().Get("date")

	if dateStr == "" {
		dateStr = time.Now().Format(time.RFC3339)
	}
	date, err := time.Parse(time.RFC3339, dateStr)

	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	weekAgo := date.AddDate(0, 0, -7)

	payments := &[]db.PurchaseSubscription{}

	query :=
		db.Database.Model(payments).
			Where("purchase_subscription.user_id = ?", c.Get("userId")).
			Where("status = ?", common.Statuses[common.APPROVED]).
			Where("purchase_subscription.created_at between ? and ? ", weekAgo, date).
			Relation("Purchase").
			Relation("Purchase.User.full_name")

	if category != "" {
		query.Where("category = ? ", category)
	}

	err = query.Select()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	count, err = db.Database.Model(payments).
		Where("purchase_subscription.user_id = ?", c.Get("userId")).
		Count()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"payments": *payments,
		"total":    count,
	}
	return c.JSON(http.StatusOK, response)
}
