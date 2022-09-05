package auth

import (
	"fmt"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/db"
	"github.com/mustafabalila/golang-api/utils/logger"
)

// ValidateSession returns the current user if the auth token is valid.
func ValidateSession(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var _, err error
	userId := fmt.Sprintf("%s", c.Get("userId"))
	var user = &db.User{Id: userId}

	err = db.Database.Model(user).WherePK().Select()

	if err == pg.ErrNoRows {
		return c.JSON(http.StatusForbidden, "Invalid token")
	}
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	response := map[string]interface{}{
		"user": user,
	}
	return c.JSON(http.StatusOK, response)
}
