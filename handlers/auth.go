package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/config"
	"github.com/mustafabalila/golang-api/db"
	"github.com/mustafabalila/golang-api/utils/logger"
	"github.com/mustafabalila/golang-api/utils/validator"
)

func createUser(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var _, err error
	auth := &CreateUser{}
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

func loginUser(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var _, err error
	var user = &db.User{}
	auth := &Login{}
	err = c.Bind(auth)
	if err != nil {
		logger.Error(err.Error())
	}

	validate := validator.New()
	err = validate.Struct(auth)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = db.Database.Model(user).Where("email = ?", auth.Email).Select()

	if err == pg.ErrNoRows {
		return c.JSON(http.StatusForbidden, "Invalid email or password")
	}

	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusForbidden, "Invalid email or password")
	}

	match, err := user.VerifyPassword(auth.Password)
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusForbidden, "Invalid email or password")
	}

	if !match {
		return c.JSON(http.StatusForbidden, "Invalid email or password")
	}

	claims := &jwt.RegisteredClaims{
		ID: user.Id,
	}

	var tokenString string
	cfg := config.GetConfig()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(cfg.JWTSecret))

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	user.FirebaseToken = auth.FirebaseToken
	db.Database.Model(user).WherePK().Update()
	response := map[string]interface{}{
		"token":  tokenString,
		"userId": user.Id,
	}
	return c.JSON(http.StatusOK, response)
}

func validateSession(c echo.Context) (e error) {
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
