package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/db"
	"github.com/mustafabalila/golang-api/utils/logger"
)

// ListUsers returns a list of users.
func ListUsers(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	users := &[]db.User{}
	err := db.Database.Model(users).Where("id != ?", c.Get("userId")).Column("id").Column("full_name").Select()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}
