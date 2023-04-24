package config

import (
	"github.com/spf13/viper"

	"github.com/turing-era/turingera-shared/log"
)

// Init 配置初始化
func Init() error {
	viper.AddConfigPath("../conf")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	log.InitLog()
	return nil
}
