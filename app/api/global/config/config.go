package config

import "time"

type Config struct {
	DatabaseConfig DatabaseConfig `mapstructure:"database" json:"database"`
	SeverConfig    SeverConfig    `mapstructure:"server" json:"sever"`
}
type SeverConfig struct {
	Mode string `mapstructure:"mode" json:"mode"`
	Host string `mapstructure:"host" json:"host"`
	Port string `mapstructure:"port" json:"port"`
}
type DatabaseConfig struct {
	MysqlConfig `mapstructure:"mysql" json:"mysql"`
	RedisConfig `mapstructure:"redis" json:"redis"`
}
type MysqlConfig struct {
	Addr            string        `mapstructure:"addr" json:"addr"`
	Port            string        `mapstructure:"port" json:"port"`
	DB              string        `mapstructure:"db" json:"db"`
	Username        string        `mapstructure:"username" json:"username"`
	Password        string        `mapstructure:"password" json:"password"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime" json:"connMaxLifetime"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns" json:"maxIdleConns"`
	ConnMaxIdleTime time.Duration `mapstructure:"connMaxIdleTime" json:"connMaxIdleTime"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns" json:"maxOpenConns"`
}
type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     string `mapstructure:"port" json:"port"`
	Username string `mapstructure:"username" json:"username"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
}
