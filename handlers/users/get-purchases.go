package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/db"
	"github.com/mustafabalila/golang-api/utils/logger"
)

// GetUserPurchases list the purchases of a user.
func GetUserPurchases(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var err error

	purchases := &[]db.Purchase{}
	err = db.Database.Model(purchases).
		Where("user_id = ?", c.Get("userId")).
		Select()

	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, purchases)
}
