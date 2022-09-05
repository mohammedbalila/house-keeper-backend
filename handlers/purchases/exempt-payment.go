package purchases

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/db"
	"github.com/mustafabalila/golang-api/notifications"
	"github.com/mustafabalila/golang-api/utils/common"
	"github.com/mustafabalila/golang-api/utils/logger"
)

// ExemptPayment exempt a purchase by updating the status of its purchase subscriptions
// to "approved" and sending a notification to the users who are subscribed
// to the purchase.
func ExemptPayment(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var err error
	purchaseId := c.Param("purchaseId")
	userId := c.Get("userId")

	purchase := &db.Purchase{Id: purchaseId}
	err = db.Database.Model(purchase).WherePK().Select()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if purchase.UserId != userId {
		return c.JSON(http.StatusUnauthorized, "You're not allowed")
	}

	payments := &[]db.PurchaseSubscription{}
	err = db.Database.Model(payments).Where("purchase_id = ?", purchaseId).Select()
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
	for _, payment := range *payments {
		payment.Status = common.Statuses[common.APPROVED]
		fmt.Printf("%s\n", payment)
		_, err = tx.Model(&payment).WherePK().Column("status").Update()
		if err != nil {
			tx.Rollback()
			logger.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
	}
	err = tx.Commit()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"payments": payments,
	}

	ids := []string{}
	for _, payment := range *payments {
		ids = append(ids, payment.UserId)
	}
	message := fmt.Sprintf("Purchase (%s) was exempted. You no longer have to pay it", purchase.Name)
	notifications.NotifyUsersWithIds(ids, message, "Purchase exempted")

	return c.JSON(http.StatusOK, response)
}
