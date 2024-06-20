package infra

import (
	"tx-monitoring/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/rs/zerolog/log"
)

var db *gorm.DB

func InitDB() {
	var err error

	db, err = gorm.Open(sqlite.Open("tx.db"), &gorm.Config{})
	if err != nil {
		log.Panic().Err(err).Msg("infra::InitDB(): failed to open db connection")
	}

	db.AutoMigrate(&model.Transaction{})
}

func GetDB() *gorm.DB {
	return db
}
