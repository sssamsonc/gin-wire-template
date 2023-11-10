package main

import (
	_ "gin-wire-template/docs"
	"gin-wire-template/routes"
	"gin-wire-template/utils/log_util"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

func main() {
	time.Local = time.UTC //make sure the timezone is UTC in this project

	logger := InitLog()
	defer logger.Sync()

	cs, err := InitApiControllers()
	if err != nil {
		log_util.Logger.Fatal("Error when initializing api controllers:" + err.Error())
	}

	r := gin.Default()
	routes.AddApi(r, cs)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run()

}
