package handlers

import (
	"net/http"
	"time"

	"github.com/go-pg/pg/v10/orm"
	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/models"
	"github.com/mustafabalila/golang-api/utils/logger"
)

func (h DBHandler) getUsers(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	users := &[]models.User{}
	err := h.DB.Model(users).Where("id != ?", c.Get("userId")).Column("id").Column("full_name").Select()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

func (h DBHandler) getUserStatistics(c echo.Context) (e error) {
	return c.JSON(http.StatusCreated, map[string]interface{}{})
}

func (h DBHandler) getUserPurchases(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var err error

	purchases := &[]models.Purchase{}
	err = h.DB.Model(purchases).
		Where("user_id = ?", c.Get("userId")).
		Where("is_complete = false").
		Select()

	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, purchases)
}

func (h DBHandler) getUserPayments(c echo.Context) (e error) {
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

	payments := &[]models.PurchaseSubscription{}

	query :=
		h.DB.Model(payments).
			Where("purchase_subscription.user_id = ?", c.Get("userId")).
			Where("status = ?", Statuses["approved"]).
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
	count, err = h.DB.Model(payments).
		Where("purchase_subscription.user_id = ?", c.Get("userId")).
		Where("status = ?", Statuses["approved"]).
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

func (h DBHandler) getUserMadePaymentRequests(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var err error

	userId := c.Get("userId")
	status := c.QueryParams().Get("status")

	requests := &[]models.PurchaseSubscription{}

	query :=
		h.DB.Model(requests).
			Where("purchase_subscription.user_id = ? ", userId).
			Where("status != ? ", Statuses["created"]).
			Relation("Purchase.id").
			Relation("Purchase.name").
			Relation("Purchase.User.full_name").
			Relation("Purchase.created_at").
			Relation("Purchase.share_price")

	if status != "" {
		query.Where("status = ? ", status)
	}

	err = query.Select()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"requests": requests,
	}
	return c.JSON(http.StatusOK, response)
}

func (h DBHandler) getOthersMadePaymentRequests(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var err error

	userId := c.Get("userId")
	status := c.QueryParams().Get("status")

	requests := &[]models.PurchaseSubscription{}

	query :=
		h.DB.Model(requests).
			Relation("Purchase.id").
			Relation("Purchase.name").
			Relation("Purchase.share_price").
			Relation("User.full_name").
			Relation("Purchase.created_at").
			Where("Purchase.user_id = ? ", userId)

	if status != "" {
		query.Where("status = ? ", status)
	}

	err = query.Select()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"requests": requests,
	}
	return c.JSON(http.StatusOK, response)
}

func (h DBHandler) getPaymentRequest(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var err error
	userId := c.Get("userId")
	id := c.Param("id")

	request := &models.PurchaseSubscription{Id: id}

	err = h.DB.Model(request).
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

func QueryToString(q *orm.Query) string {
	value, _ := q.AppendQuery(orm.NewFormatter(), nil)

	return string(value)
}
