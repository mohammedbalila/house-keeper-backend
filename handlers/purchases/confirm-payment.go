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

// ConfirmPayment confirms a payment (purchase subscription) by updating its
// status to "approved" and notifies the payment creator (purchase subscription user).
func ConfirmPayment(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var err error

	subscription := &db.PurchaseSubscription{Id: c.Param("purchaseSubscriptionId")}
	err = db.Database.Model(subscription).WherePK().Relation("Purchase.user_id").Select()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if subscription.Purchase.UserId != c.Get("userId") {
		return c.JSON(http.StatusUnauthorized, "You're not allowed")
	}

	payment := &db.PurchaseSubscription{
		Id:     subscription.Id,
		Status: common.Statuses[common.APPROVED],
	}

	_, err = db.Database.Model(payment).WherePK().Column("status").Update()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"payment": payment,
	}

	message := fmt.Sprintf("Your payment to %s was approved by %s. Thanks for your cooperation",
		subscription.Purchase.Name,
		subscription.Purchase.User.FullName)
	notifications.NotifyUserByToken(notifications.NotifyInput{Token: subscription.User.FirebaseToken, Title: "Payment approved", Body: message})

	return c.JSON(http.StatusOK, response)
}
