package http_util

import (
	"gin-wire-template/configs/common_config"
	"net/http"
)

func NewClient() *http.Client {
	commonConfig := common_config.NewConfig()

	return &http.Client{
		Timeout: commonConfig.HttpClientTimeOut,
	}
}
