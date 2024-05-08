package config

import "os"

type Database struct {
	Localhost string
	User      string
	Password  string
	Port      string
	DbName    string
	SSLMode   string
	TimeZone  string
}

type Mail struct {
	Host     string
	Port     string
	Username string
	Password string
}

type Service struct {
	AccountURL string
}

type Config struct {
	Database Database
	Mail     Mail
	Service  Service
}

func Load() *Config {
	return &Config{
		Database: Database{
			Localhost: os.Getenv("POSTGRES_HOST"),
			User:      os.Getenv("POSTGRES_USER"),
			Password:  os.Getenv("POSTGRES_PASSWORD"),
			Port:      os.Getenv("POSTGRES_PORT"),
			DbName:    os.Getenv("POSTGRES_DB"),
			SSLMode:   os.Getenv("POSTGRES_SLLMODE"),
			TimeZone:  os.Getenv("POSTGRES_TIMEZONE"),
		},
		Mail: Mail{
			Host:     os.Getenv("MAIL_HOST"),
			Port:     os.Getenv("MAIL_PORT"),
			Username: os.Getenv("MAIL_USERNAME"),
			Password: os.Getenv("MAIL_PASSWORD"),
		},
		Service: Service{
			AccountURL: os.Getenv("ACCOUNT_URL"),
		},
	}
}
