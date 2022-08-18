package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/config"
)

func RegisterRoutes(prefix *echo.Group, h DBHandler) {

	auth := prefix.Group("/auth")
	auth.POST("/register", h.createUser)
	auth.POST("/login", h.loginUser)
	auth.GET("/session", h.validateSession, config.AuthMiddleware())

	users := prefix.Group("/users")
	users.Use(config.AuthMiddleware())
	users.GET("/", h.getUsers)
	users.GET("/purchases", h.getUserPurchases)
	users.GET("/payments", h.getUserPayments)
	users.GET("/my-requests", h.getUserMadePaymentRequests)
	users.GET("/others-requests", h.getOthersMadePaymentRequests)
	users.GET("/requests/:id", h.getPaymentRequest)
	users.GET("/statistics", h.getUserStatistics)

	purchases := prefix.Group("/purchases")

	purchases.Use(config.AuthMiddleware())
	purchases.POST("/", h.createPurchase)
	purchases.GET("/", h.getUnPaidPurchases)
	purchases.GET("/:purchaseId", h.getPurchaseDetail)
	purchases.POST("/request-confirmation/:purchaseId", h.requestPaymentConformation)
	purchases.POST("/confirm/:purchaseSubscriptionId", h.confirmPayment)
	purchases.POST("/reject/:purchaseSubscriptionId", h.rejectPayment)
	purchases.POST("/exempt/:purchaseId", h.exemptPayment)
}
