package database

import (
	log "github.com/sirupsen/logrus"
	"transact-api/configuration"

	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db    *gorm.DB
	once  sync.Once
	dbErr error
)

func ConnectDB(config configuration.Configuration) (*gorm.DB, error) {
	once.Do(func() {
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			config.GetDBHost(),
			config.GetDBUser(),
			config.GetDBPassword(),
			config.GetDBName(),
			config.GetDBPort(),
		)

		log.Println("dsn", dsn)
		db, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	})

	return db, dbErr
}
