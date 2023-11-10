package databases

import (
	"gin-wire-template/databases/mongo_database"
	"gin-wire-template/databases/mysql_database"
	"gin-wire-template/databases/redis_cache"
	"github.com/google/wire"
)

var DatabaseSet = wire.NewSet(
	mongo_database.NewConnector,
	mysql_database.NewConnector,
	redis_cache.NewConnector,
)
