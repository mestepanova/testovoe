package app

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort string
	AppHost string

	DBPort     string
	DBHost     string
	DBName     string
	DBUser     string
	DBPassword string
}

func InitConfig() *Config {
	return &Config{
		AppPort:    os.Getenv("APP_PORT"),
		AppHost:    os.Getenv("APP_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBName:     os.Getenv("DB_NAME"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}
}

func (c *Config) DSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=disable"
}

func (c *Config) MigrateDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName,
	)
}

func (c *Config) Addr() string {
	return c.AppHost + ":" + c.AppPort
}
