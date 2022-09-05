package purchases

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/db"
	"github.com/mustafabalila/golang-api/utils/logger"
)

// GetPurchase retrieve a purchase with its subscriptions count by the purchase Id.
// It also retrieves its purchase subscriptions if the authenticated user
// is the purchase owner.
func GetPurchase(c echo.Context) (e error) {
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
