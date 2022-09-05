package auth

import (
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/db"
	"github.com/mustafabalila/golang-api/utils/logger"
	"github.com/mustafabalila/golang-api/utils/validator"
)

type createUserInput struct {
	FullName string `json:"fullName" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
}

// SignUp creates a new user and returns it.
func SignUp(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var _, err error
	auth := &createUserInput{}
	err = c.Bind(&auth)
	if err != nil {
		return err
	}

	validate := validator.New()
	err = validate.Struct(auth)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	user := &db.User{
		FullName: auth.FullName,
		Password: auth.Password,
		Email:    auth.Email,
	}

	// check if user already exists
	existingUser := &db.User{}
	err = db.Database.Model(existingUser).Where("email = ?", user.Email).Select()
	if err != pg.ErrNoRows {
		return c.JSON(http.StatusConflict, "User already exists")
	}

	err = user.HashPassword()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	_, err = db.Database.Model(user).Insert()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"user": user,
	}
	return c.JSON(http.StatusCreated, response)
}
