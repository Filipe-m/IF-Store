package platform

import (
	"account/internal/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnect(cfg config.Database) (*gorm.DB, error) {

	postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Localhost, cfg.User, cfg.Password, cfg.DbName, cfg.Port, cfg.SSLMode, cfg.TimeZone)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  postgresDSN,
		PreferSimpleProtocol: true,
	}))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *gorm.DB, entities ...interface{}) error {
	for _, entity := range entities {
		if err := db.AutoMigrate(entity); err != nil {
			return err
		}
	}
	return nil
}
