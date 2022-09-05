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

// RequestPaymentConformation requests a payment confirmation from the user by updating
// the status of the purchase subscription to "pending" and sending a notification
// to the user who owns the purchase.
func RequestPaymentConformation(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var err error
	userId := fmt.Sprintf("%s", c.Get("userId"))
	purchaseId := c.Param("purchaseId")

	payment := &db.PurchaseSubscription{
		PurchaseId: purchaseId,
		UserId:     userId,
	}

	err = db.Database.Model(payment).
		Where("purchase_id = ?", purchaseId).
		Where("purchase_subscription.user_id = ? ", userId).
		Relation("User.full_name").
		Relation("Purchase.name").
		Relation("Purchase.User.firebase_token").
		Select()

	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// can only confirm payment if the status is created
	if payment.Status != common.Statuses[common.CREATED] {
		return c.JSON(http.StatusBadRequest, "Can't do")
	}

	payment.Status = common.Statuses[common.PENDING]
	_, err = db.Database.Model(payment).Where("user_id = ?", userId).Where("purchase_id = ?", purchaseId).Column("status").Update()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"payment": payment,
	}

	message := fmt.Sprintf("You have a new payment request on %s by %s.",
		payment.Purchase.Name,
		payment.User.FullName)
	notifications.NotifyUserByToken(notifications.NotifyInput{Token: payment.Purchase.User.FirebaseToken, Title: "Payment approval request", Body: message})

	return c.JSON(http.StatusOK, response)
}
