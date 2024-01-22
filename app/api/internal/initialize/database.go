package initialize

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"gobang/app/api/global"
	"time"
)

func SetupDatabase() {
	setupMysql()
	setupRedis()
}
func setupMysql() {
	config := global.ConfigName.DatabaseConfig.MysqlConfig
	//dsn := "username:password@tcp(127.0.0.1:3306)/database_name?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := config.Username + ":" + config.Password + "@tcp(" + config.Addr + ":" + config.Port + ")/" + config.DB + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		global.Logger.Error("open mysql failed" + err.Error())
		return
	}
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	db.SetMaxOpenConns(config.MaxOpenConns)
	err = db.Ping()
	if err != nil {
		global.Logger.Fatal("connect to my sql failed:" + err.Error())
	}
	global.MysqlDB = db
	global.Logger.Info("init mysql success")
}
func setupRedis() {
	config := global.ConfigName.DatabaseConfig.RedisConfig
	rdb := redis.NewClient(&redis.Options{
		Addr: config.Host + ":" + config.Port,
		//ClientName:            "",
		Username: config.Username,
		Password: config.Password,
		DB:       config.DB,
	})
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		global.Logger.Fatal("connect to redis failed" + err.Error())
	}
	global.RDB = rdb
	global.Logger.Info("init redis success")
}
