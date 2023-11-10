//go:build wireinject
// +build wireinject

package main

import (
	"gin-wire-template/configs"
	"gin-wire-template/configs/common_config"
	"gin-wire-template/controllers"
	"gin-wire-template/databases"
	"gin-wire-template/repositories"
	"gin-wire-template/utils/log_util"
	"github.com/google/wire"
	"go.uber.org/zap"
)

func InitApiControllers() (*controllers.Controllers, error) {
	wire.Build(
		//configs
		configs.ConfigSet,
		//databases
		databases.DatabaseSet,
		//controllers
		controllers.ControllerSet,
		//repositories
		repositories.RepositorySet,
	)
	return nil, nil
}

func InitLog() *zap.Logger {
	wire.Build(
		log_util.NewLogger,
		common_config.NewConfig,
	)
	return nil
}
