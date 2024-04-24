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

type Service struct {
	InventoryUrl    string
	NotificationUrl string
	PaymentUrl      string
	ShipmentUrl     string
}

type Config struct {
	Database Database
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
		Service: Service{
			InventoryUrl:    os.Getenv("INVENTORY_URL"),
			NotificationUrl: os.Getenv("NOTIFICATION_URL"),
			PaymentUrl:      os.Getenv("PAYMENT_URL"),
			ShipmentUrl:     os.Getenv("SHIPMENT_URL"),
		},
	}
}
