package db

import (
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/mustafabalila/golang-api/config"
	"github.com/mustafabalila/golang-api/utils/logger"
	"go.uber.org/zap"
)

var Database *pg.DB

// Connect to the database and sets the value of the Database variable
func Connect() {
	log := logger.GetLoggerInstance()
	defer log.Sync()

	var err error
	cfg := config.GetConfig()
	opt, err := pg.ParseURL(cfg.DatabaseUrl)
	if err != nil {
		log.Fatal("Failed to parse database url", zap.Error(err))
		os.Exit(1)
	}
	Database = pg.Connect(opt)
}
