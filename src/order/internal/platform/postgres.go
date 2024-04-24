package platform

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"order/internal/config"
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

func CommitOrRollback(tx *gorm.DB, err error) error {
	if tx == nil {
		return fmt.Errorf("no transaction to commit")
	}

	if err != nil {
		if err = tx.Rollback().Error; err != nil {
			return err
		}
		return err
	}
	return tx.Commit().Error
}
