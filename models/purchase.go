package models

import "fmt"

type Purchase struct {
	Base
	Id              string  `json:"id"`
	UserId          string  `json:"userId"`
	Name            string  `json:"name"`
	TotalPrice      float64 `json:"totalPrice"`
	SharePrice      float64 `json:"sharePrice"`
	Description     string  `json:"description"`
	PaymentProgress int     `json:"paymentProgress"`
	User            *User   `json:"user" pg:"rel:has-one"`
	Category        int     `json:"category"`

	tableName struct{} `pg:"api.purchase"`
}

type PurchaseSubscription struct {
	Base
	Id         string    `json:"id"`
	UserId     string    `json:"userId"`
	PurchaseId string    `json:"purchaseId"`
	Status     int64     `json:"status"`
	Purchase   *Purchase `json:"purchase"`
	User       *User     `json:"user"`
	tableName  struct{}  `pg:"api.purchase_subscription"`
}

func (t Purchase) String() string {
	return fmt.Sprintf("Purchase<%s %s %f>", t.Id, t.Name, t.TotalPrice)
}

func (t PurchaseSubscription) String() string {
	return fmt.Sprintf("PurchaseSubscription<%s %s %s>", t.Id, t.UserId, t.PurchaseId)
}
