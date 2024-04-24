package platform

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnect(cfg Database) (*gorm.DB, error) {

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

type Database struct {
	Localhost string
	User      string
	Password  string
	Port      string
	DbName    string
	SSLMode   string
	TimeZone  string
}
