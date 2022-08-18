package config

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"

	"github.com/labstack/echo/v4"
)

// TokenParser represents JWT token parser
type TokenParser interface {
	ParseToken(string) (*jwt.Token, error)
}

// Middleware makes JWT implement the Middleware interface.
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(header, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(GetConfig().JWTSecret), nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, err.Error())
			}

			id := claims["jti"]
			c.Set("userId", id)

			return next(c)
		}
	}
}
