package auth

import (
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/config"
	"github.com/mustafabalila/golang-api/db"
	"github.com/mustafabalila/golang-api/utils/logger"
	"github.com/mustafabalila/golang-api/utils/validator"
)

type loginInput struct {
	Password      string `json:"password" validate:"required,min=8"`
	Email         string `json:"email" validate:"required,email"`
	FirebaseToken string `json:"firebaseToken" validate:"required"`
}

// Login logs in a user and returns a JWT token and the user Id.
func Login(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var _, err error
	var user = &db.User{}
	auth := &loginInput{}
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
