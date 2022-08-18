package handlers

import "github.com/go-pg/pg/v10"

type DBHandler struct {
	DB pg.DB
}
