package db

import (
	"fmt"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/mustafabalila/golang-api/config"
	"github.com/mustafabalila/golang-api/utils/logger"
	"go.uber.org/zap"
)

var Database *pg.DB

type Attachment struct {
	Id               string
	BucketName       string
	AttachmentPath   string
	OriginalFileName string
	MimeType         string
	Size             int64

	tableName struct{} `pg:"api.attachment"`
}

type Page struct {
	Id           string
	BookId       string
	ChapterId    string
	Index        int
	AttachmentId string

	tableName struct{} `pg:"api.page"`
}

func Connect() {
	log := logger.GetLoggerInstance()
	defer log.Sync()
	var err error = nil
	cfg := config.GetConfig()
	fmt.Println(cfg.DatabaseUrl)
	opt, err := pg.ParseURL(cfg.DatabaseUrl)
	if err != nil {
		log.Fatal("Failed to parse database url", zap.Error(err))
		os.Exit(1)
	}
	Database = pg.Connect(opt)
}
