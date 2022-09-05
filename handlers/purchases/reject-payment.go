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

// GetUnPaidPurchases reject a payment (purchase subscription) by updating its status
// to "rejected" and sending a notification to the user who made the payment.
func RejectPayment(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var err error

	subscription := &db.PurchaseSubscription{Id: c.Param("purchaseSubscriptionId")}
	err = db.Database.Model(subscription).
		WherePK().
		Relation("User.firebase_token").
		Relation("Purchase.name").
		Relation("Purchase.user_id").
		Relation("Purchase.User.full_name").
		Select()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if subscription.Purchase.UserId != c.Get("userId") {
		return c.JSON(http.StatusUnauthorized, "You're not allowed")
	}

	// can only reject payment if the status is pending
	if subscription.Status != common.Statuses[common.PENDING] {
		return c.JSON(http.StatusBadRequest, "Can't do")
	}

	payment := &db.PurchaseSubscription{
		Id:     subscription.Id,
		Status: common.Statuses[common.REJECTED],
	}

	_, err = db.Database.Model(payment).WherePK().Column("status").Update()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"payment": payment,
	}

	message := fmt.Sprintf("Your payment to %s was rejected by %s. Please refer to them for more details",
		subscription.Purchase.Name,
		subscription.Purchase.User.FullName)
	notifications.NotifyUserByToken(notifications.NotifyInput{Token: subscription.User.FirebaseToken, Title: "Payment rejected", Body: message})
	return c.JSON(http.StatusOK, response)
}
