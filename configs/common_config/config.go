package common_config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	IsDemoMode        bool
	ShowDebugLog      bool
	HttpClientTimeOut time.Duration
}

func NewConfig() (config Config) {
	httpClientTimeOut := time.Duration(30) * time.Second //default
	v, isExist := os.LookupEnv("HTTP_CLIENT_TIME_OUT_SECONDS")
	if isExist {
		x, err := strconv.Atoi(v)
		if err == nil {
			httpClientTimeOut = time.Duration(x) * time.Second
		}
	}

	return Config{
		IsDemoMode:        os.Getenv("IS_DEMO_MODE") == "true",
		ShowDebugLog:      os.Getenv("SHOW_DEBUG_LOG") == "true",
		HttpClientTimeOut: httpClientTimeOut,
	}
}
