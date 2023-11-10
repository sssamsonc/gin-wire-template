package mysql_db_config

import (
	"errors"
	"os"
)

type Config struct {
	Host     string
	Database string
	Username string
	Password string
}

func NewConfig() (Config, error) {
	if os.Getenv("MYSQL_DB_HOST") == "" {
		return Config{}, errors.New("missing mandatory environment variables: MYSQL_DB_HOST")
	}

	return Config{
		Host:     os.Getenv("MYSQL_DB_HOST"),
		Database: os.Getenv("MYSQL_DB_DATABASE"),
		Username: os.Getenv("MYSQL_DB_USERNAME"),
		Password: os.Getenv("MYSQL_DB_PASSWORD"),
	}, nil
}
