package database

import (
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

		db, dbErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	})

	return db, dbErr
}
