package global

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gobang/app/api/global/config"
	"gobang/app/api/internal/model"
)

var (
	ConfigName config.Config
	Logger     *zap.Logger
	MysqlDB    *sqlx.DB
	RDB        *redis.Client
	Counter    model.Counter
)
