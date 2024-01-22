package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gobang/app/api/global"
	"time"
)

func CheckUser(username string) (bool, error) {
	ctx := context.Background()
	added, err := global.RDB.SAdd(ctx, "user_ids", username).Result()
	if err != nil {
		global.Logger.Warn("redis CheckUser failed" + err.Error())
		return true, err
	}
	if added == 0 { //如果存在重复的名字
		global.Logger.Info("username " + username + " existed in redis")
		return true, err
	}
	return false, err
}
func AddUser(username string, password string) error {

	ctx := context.Background()
	err := global.RDB.HSet(ctx, fmt.Sprintf("username:%s", username), "password", password).Err()
	if err != nil {
		global.Logger.Warn("Add user failed in redis" + err.Error())
		return err
	}
	// 设置过期时间
	err = global.RDB.Expire(ctx, fmt.Sprintf("user:%s", username), 3*time.Hour).Err()
	if err != nil {
		global.Logger.Warn("set user expire time failed in redis" + err.Error())
		return err
	}
	return err
}
func CheckPassword(username string, password string) (bool, error) {
	hashKey := fmt.Sprintf("user:%s", username)
	ctx := context.Background()
	GetPassword, err := global.RDB.HGet(ctx, hashKey, "password").Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		global.Logger.Warn("get password failed in redis:" + err.Error())
		return false, err
	}
	if GetPassword == password {
		return true, err
	}
	return false, err
}
