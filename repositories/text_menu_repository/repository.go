package text_menu_repository

import (
	"gin-wire-template/configs/common_config"
	"gin-wire-template/configs/redis_config"
	"gin-wire-template/databases/mongo_database"
	"gin-wire-template/databases/mysql_database"
	"gin-wire-template/databases/redis_cache"
)

const (
	COLLECTION_TEXT_MENU = "text_menu"
	TEXT_MENU_CACHE_NAME = "text_menu_cache"
)

type Repository struct {
	commonConfig   common_config.Config
	mongoConnector *mongo_database.Connector
	mysqlConnector *mysql_database.Connector
	redisConfig    redis_config.Config
	redisConnector *redis_cache.Connector
}

func NewRepository(
	commonConfig common_config.Config,
	mongoConnector *mongo_database.Connector,
	mysqlConnector *mysql_database.Connector,
	redisConfig redis_config.Config,
	redisConnector *redis_cache.Connector,
) *Repository {
	return &Repository{
		commonConfig: commonConfig,
		//mongo
		mongoConnector: mongoConnector,
		//mysql
		mysqlConnector: mysqlConnector,
		//redis
		redisConfig:    redisConfig,
		redisConnector: redisConnector,
	}
}
