package initialize

import (
	"github.com/spf13/viper"
	"gobang/app/api/global"
)

func SetupViper() {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.SetConfigFile("./manifest/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic("viper read config failed" + err.Error())
	}
	err = viper.Unmarshal(&global.ConfigName)
	if err != nil {
		panic("viper unmarshal failed" + err.Error())
	}
}
