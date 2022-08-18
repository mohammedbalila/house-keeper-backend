package handlers

type CreateUser struct {
	FullName string `json:"fullName" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
}

type Login struct {
	Password      string `json:"password" validate:"required,min=8"`
	Email         string `json:"email" validate:"required,email"`
	FirebaseToken string `json:"firebaseToken" validate:"required"`
}

type CreatePurchase struct {
	Name            string   `json:"name"`
	TotalPrice      float64  `json:"totalPrice"`
	Description     string   `json:"description"`
	Category        int      `json:"category"`
	PaymentProgress int      `json:"paymentProgress"`
	Subscribers     []string `json:"subscribers"`
}

var Statuses = map[string]int64{
	"created":  1,
	"pending":  2,
	"approved": 3,
	"rejected": 4,
}
