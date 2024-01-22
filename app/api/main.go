package main

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gobang/app/api/global"
	"gobang/app/api/internal/initialize"
	"gobang/app/api/router"
)

func main() {
	initialize.SetupViper()
	initialize.SetupLogger()
	initialize.SetupDatabase()
	config := global.ConfigName.SeverConfig
	gin.SetMode(config.Mode)

	global.Logger.Info("init server success", zap.String("port", config.Port+":"+config.Port))
	err := router.InitRouter(config.Port)
	if err != nil {
		global.Logger.Fatal("server start failed" + err.Error())
	}

}
