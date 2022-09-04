package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/config"
	"github.com/mustafabalila/golang-api/db"
)

func RegisterRoutes(prefix *echo.Group) {

	db.Connect()
	auth := prefix.Group("/auth")
	auth.POST("/register", createUser)
	auth.POST("/login", loginUser)
	auth.GET("/session", validateSession, config.AuthMiddleware())

	users := prefix.Group("/users")
	users.Use(config.AuthMiddleware())
	users.GET("/", getUsers)
	users.GET("/purchases", getUserPurchases)
	users.GET("/payments", getUserPayments)
	users.GET("/my-requests", getUserMadePaymentRequests)
	users.GET("/others-requests", getOthersMadePaymentRequests)
	users.GET("/requests/:id", getPaymentRequest)
	users.GET("/statistics", getUserStatistics)

	purchases := prefix.Group("/purchases")

	purchases.Use(config.AuthMiddleware())
	purchases.POST("/", createPurchase)
	purchases.GET("/", getUnPaidPurchases)
	purchases.GET("/:purchaseId", getPurchaseDetail)
	purchases.POST("/request-confirmation/:purchaseId", requestPaymentConformation)
	purchases.POST("/confirm/:purchaseSubscriptionId", confirmPayment)
	purchases.POST("/reject/:purchaseSubscriptionId", rejectPayment)
	purchases.POST("/exempt/:purchaseId", exemptPayment)
}
