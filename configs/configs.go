package configs

import (
	"gin-wire-template/configs/common_config"
	"gin-wire-template/configs/mongo_db_config"
	"gin-wire-template/configs/mysql_db_config"
	"gin-wire-template/configs/redis_config"
	"github.com/google/wire"
)

var ConfigSet = wire.NewSet(
	common_config.NewConfig,
	mongo_db_config.NewConfig,
	mysql_db_config.NewConfig,
	redis_config.NewConfig,
)
