package router

import (
	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/config"
	"github.com/mustafabalila/golang-api/db"
	"github.com/mustafabalila/golang-api/handlers/auth"
	"github.com/mustafabalila/golang-api/handlers/purchases"
	"github.com/mustafabalila/golang-api/handlers/users"
)

// SetupRoutes sets up the routes for the server
func SetupRoutes(e *echo.Echo) {
	api := e.Group("/api")
	prefix := api.Group("/v1")
	db.Connect()

	authRoutes := prefix.Group("/auth")
	authRoutes.POST("/register", auth.SignUp)
	authRoutes.POST("/login", auth.Login)
	authRoutes.GET("/session", auth.ValidateSession, config.AuthMiddleware())

	userRoutes := prefix.Group("/users")
	userRoutes.Use(config.AuthMiddleware())
	userRoutes.GET("/", users.ListUsers)
	userRoutes.GET("/purchases", users.GetUserPurchases)
	userRoutes.GET("/payments", users.GetUserPayments)
	userRoutes.GET("/my-requests", users.GetSelfPaymentRequests)
	userRoutes.GET("/others-requests", users.GetOthersPaymentRequests)
	userRoutes.GET("/requests/:id", users.RetrievePaymentRequest)

	purchaseRoutes := prefix.Group("/purchases")
	purchaseRoutes.Use(config.AuthMiddleware())
	purchaseRoutes.POST("/", purchases.CreatePurchase)
	purchaseRoutes.GET("/", purchases.GetUnPaidPurchases)
	purchaseRoutes.GET("/:purchaseId", purchases.GetPurchase)
	purchaseRoutes.POST("/request-confirmation/:purchaseId", purchases.RequestPaymentConformation)
	purchaseRoutes.POST("/confirm/:purchaseSubscriptionId", purchases.ConfirmPayment)
	purchaseRoutes.POST("/reject/:purchaseSubscriptionId", purchases.RejectPayment)
	purchaseRoutes.POST("/exempt/:purchaseId", purchases.ExemptPayment)
}
