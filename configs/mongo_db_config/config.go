package mongo_db_config

import (
	"errors"
	"os"
)

type Config struct {
	Host       string
	Database   string
	Username   string
	Password   string
	Replicaset string
}

func NewConfig() (Config, error) {
	if os.Getenv("MONGO_DB_HOST") == "" {
		return Config{}, errors.New("missing mandatory environment variables: MONGO_DB_HOST")
	}

	return Config{
		Host:       os.Getenv("MONGO_DB_HOST"),
		Database:   os.Getenv("MONGO_DB_DATABASE"),
		Username:   os.Getenv("MONGO_DB_USERNAME"),
		Password:   os.Getenv("MONGO_DB_PASSWORD"),
		Replicaset: os.Getenv("MONGO_DB_REPLICASET"),
	}, nil
}
