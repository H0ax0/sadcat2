package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	Token       string
	LogBasePath string
	Sql         *SqlConfig
)

type SqlConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DB       string
}

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetDefault("bot_token", "1937425133:AAFfZHpmBaCEXFkqkXq0KpeDg0mE8hGPtK8")
	viper.SetDefault("logs_path", ".")
	viper.SetDefault("sql.host", "localhost")
	viper.SetDefault("sql.port", "5432")
	viper.SetDefault("sql.user", "user")
	viper.SetDefault("sql.password", "password")
	viper.SetDefault("sql.database", "sadcat2")

	err := viper.ReadInConfig()
	if err != nil {
		zap.S().Errorw("failed to read config", "error", err)
	}

	Token = viper.GetString("bot_token")
	LogBasePath = viper.GetString("logs_path")

	Sql = &SqlConfig{
		Host:     viper.GetString("sql.host"),
		Port:     viper.GetInt("sql.port"),
		User:     viper.GetString("sql.user"),
		Password: viper.GetString("sql.password"),
		DB:       viper.GetString("sql.database"),
	}

}
