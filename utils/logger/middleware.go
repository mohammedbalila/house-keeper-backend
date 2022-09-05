package logger

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LoggerMiddleware is a middleware that logs the request information
func LoggerMiddleware() echo.MiddlewareFunc {
	log := GetLoggerInstance()
	defer log.Sync()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()

			start := time.Now()
			next(c)
			stop := time.Now()

			reqId := req.Header.Get(echo.HeaderXRequestID)
			if reqId == "" {
				reqId = res.Header().Get(echo.HeaderXRequestID)
			}

			method := c.Request().Method
			path := c.Path()
			status := c.Response().Status
			ip := c.RealIP()
			hostname := c.Request().Host

			log.Info(
				"request",
				zap.String("requestId", reqId),
				zap.String("method", method),
				zap.String("path", path),
				zap.Int("status", status),
				zap.String("ip", ip),
				zap.String("hostname", hostname),
				zap.String("latency", stop.Sub(start).String()),
			)

			return nil
		}
	}
}
