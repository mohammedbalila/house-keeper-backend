package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/db"
	"github.com/mustafabalila/golang-api/utils/logger"
	"github.com/mustafabalila/golang-api/utils/notifications"
)

func createPurchase(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var _, err error
	input := &CreatePurchase{}
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
			Status:     Statuses["created"],
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
	notifyUsers(input.Subscribers, message, "New Purchase")
	response := map[string]interface{}{
		"purchase": purchase,
	}
	return c.JSON(http.StatusCreated, response)
}

func getUnPaidPurchases(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var err error
	users := c.QueryParams().Get("users")
	createdAt := c.QueryParams().Get("createdAt")
	subscriptions := &[]db.PurchaseSubscription{}
	userId := c.Get("userId")

	query := db.Database.Model(subscriptions).
		Where("purchase_subscription.user_id = ?", userId).
		Where("status = ?", Statuses["created"]).
		Relation("Purchase").
		Relation("Purchase.User.full_name")

	if createdAt != "" {
		query.Where("purchase_subscription.created_at >= ?", createdAt)
	}
	if users != "" {
		userIds := strings.Split(users, ",")
		query.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q = q.
				Where("purchase.user_id in (?) ", pg.In(userIds)).
				Where("purchase_subscription.user_id = ?", c.Get("userId"))
			return q, nil
		})

	}
	err = query.Select()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"subscriptions": subscriptions,
	}
	return c.JSON(http.StatusOK, response)
}

func getPurchaseDetail(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var err error

	purchase := &db.Purchase{Id: c.Param("purchaseId")}
	err = db.Database.Model(purchase).WherePK().Relation("User.full_name").Select()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	if purchase.UserId == c.Get("userId") {
		payments := &[]db.PurchaseSubscription{}
		err = db.Database.Model(payments).Where("purchase_id = ? ", c.Param("purchaseId")).Relation("User.full_name").Select()
		if err != nil {
			logger.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		response := map[string]interface{}{
			"purchase":    purchase,
			"payments":    payments,
			"subscribers": len(*payments),
		}
		return c.JSON(http.StatusOK, response)
	}

	countResult := map[string]interface{}{
		"subscribers": 0,
	}

	err = db.Database.Model(&db.PurchaseSubscription{}).
		Where("purchase_id = ? ", c.Param("purchaseId")).
		ColumnExpr("count(*) as subscribers").
		Select(&countResult)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Printf("\n%s\n", countResult)
	response := map[string]interface{}{
		"purchase":    purchase,
		"payments":    []string{},
		"subscribers": countResult["subscribers"],
	}
	return c.JSON(http.StatusOK, response)
}

func requestPaymentConformation(c echo.Context) (e error) {
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

	if payment.Status != Statuses["created"] {
		return c.JSON(http.StatusBadRequest, "Can't do")
	}

	payment.Status = Statuses["pending"]
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
	notifications.NotifyUser(notifications.NotifyInput{Token: payment.Purchase.User.FirebaseToken, Title: "Payment approval request", Body: message})

	return c.JSON(http.StatusOK, response)
}

func confirmPayment(c echo.Context) (e error) {
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
		Status: Statuses["approved"],
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
	notifications.NotifyUser(notifications.NotifyInput{Token: subscription.User.FirebaseToken, Title: "Payment approved", Body: message})

	return c.JSON(http.StatusOK, response)
}

func rejectPayment(c echo.Context) (e error) {
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

	payment := &db.PurchaseSubscription{
		Id:     subscription.Id,
		Status: Statuses["rejected"],
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
	notifications.NotifyUser(notifications.NotifyInput{Token: subscription.User.FirebaseToken, Title: "Payment rejected", Body: message})
	return c.JSON(http.StatusOK, response)
}

func exemptPayment(c echo.Context) (e error) {
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
		payment.Status = Statuses["approved"]
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
	notifyUsers(ids, message, "Purchase exempted")

	return c.JSON(http.StatusOK, response)
}

func notifyUsers(ids []string, message string, title string) error {
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

	err = notifications.NotifyUsers(
		tokens,
		notifications.NotifyInput{Title: title, Body: message})

	if err != nil {
		return err
	}

	return nil
}
