package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/db"
	"github.com/mustafabalila/golang-api/utils/logger"
)

// GetSelfPaymentRequests list all payment requests (purchase subscriptions)
// that were made by the authenticated user.
// it can be filtered by status and created_at (cursor).
func GetSelfPaymentRequests(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var err error

	userId := c.Get("userId")
	status := c.QueryParams().Get("status")
	created_at := c.QueryParams().Get("created_at")

	requests := &[]db.PurchaseSubscription{}

	query :=
		db.Database.Model(requests).
			Where("purchase_subscription.user_id = ? ", userId).
			Relation("Purchase.id").
			Relation("Purchase.name").
			Relation("Purchase.User.full_name").
			Relation("Purchase.created_at").
			Relation("Purchase.share_price")

	if status != "" {
		query.Where("status = ? ", status)
	}
	if created_at != "" {
		query.Where("purchase_subscription.created_at >= ? ", created_at)
	}

	err = query.
		Order("created_at desc").
		Limit(10). // it's paginated using the created_at field (cursor)
		Select()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"requests": requests,
	}
	return c.JSON(http.StatusOK, response)
}

// GetOthersPaymentRequests lists all payment requests (purchase subscriptions)
// that were made by other users for the authenticated user purchases.
// it can be filtered by status and created_at (cursor).
func GetOthersPaymentRequests(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var err error

	userId := c.Get("userId")
	status := c.QueryParams().Get("status")
	created_at := c.QueryParams().Get("created_at")

	requests := &[]db.PurchaseSubscription{}

	query :=
		db.Database.Model(requests).
			Relation("Purchase.id").
			Relation("Purchase.name").
			Relation("Purchase.share_price").
			Relation("User.full_name").
			Relation("Purchase.created_at").
			Where("Purchase.user_id = ? ", userId)

	if status != "" {
		query.Where("status = ? ", status)
	}

	if created_at != "" {
		query.Where("purchase_subscription.created_at >= ? ", created_at)
	}

	err = query.
		Order("created_at desc").
		Limit(10). // it's paginated using the created_at field (cursor)
		Select()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"requests": requests,
	}
	return c.JSON(http.StatusOK, response)
}

// RetrievePaymentRequest retrieve a payment request (purchase subscription) by id.
func RetrievePaymentRequest(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var err error
	userId := c.Get("userId")
	id := c.Param("id")

	request := &db.PurchaseSubscription{Id: id}

	err = db.Database.Model(request).
		Relation("Purchase.id").
		Relation("Purchase.name").
		Relation("Purchase.user_id").
		Relation("Purchase.User.full_name").
		Relation("Purchase.share_price").
		Relation("User.full_name").
		Relation("Purchase.created_at").
		WherePK().
		Select()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if request.UserId != userId || request.Purchase.UserId != userId {
		return c.JSON(http.StatusForbidden, "You're not allowed")
	}
	response := map[string]interface{}{
		"request": request,
	}
	return c.JSON(http.StatusOK, response)
}
