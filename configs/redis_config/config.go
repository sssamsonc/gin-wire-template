package redis_config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Host               string
	Database           string
	Username           string
	Password           string
	CacheRetryLockTime time.Duration
	CacheTTL           time.Duration
}

func NewConfig() (Config, error) {
	cacheRetryLockSec := time.Duration(60) * time.Second //default 60 seconds
	v, isExist := os.LookupEnv("CACHE_RETRY_LOCK_SEC")
	if isExist {
		x, err := strconv.Atoi(v)
		if err == nil {
			cacheRetryLockSec = time.Duration(x) * time.Second
		}
	}

	cacheTTL := time.Duration(300) * time.Second //default 300 seconds
	v, isExist = os.LookupEnv("CACHE_TTL_SEC")
	if isExist {
		x, err := strconv.Atoi(v)
		if err == nil {
			cacheTTL = time.Duration(x) * time.Second
		}
	}

	return Config{
		Host:               os.Getenv("REDIS_HOST"),
		Database:           os.Getenv("REDIS_DATABASE"),
		Username:           os.Getenv("REDIS_USERNAME"),
		Password:           os.Getenv("REDIS_PASSWORD"),
		CacheRetryLockTime: cacheRetryLockSec,
		CacheTTL:           cacheTTL,
	}, nil
}
