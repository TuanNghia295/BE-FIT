package config

import (
	"strings"
	"sync"

	"github.com/spf13/viper" // load config
)

type Config struct {
	Server *Server
	Db     *Db
}

type Server struct {
	Port int
}

type Db struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

// var for package zone, it can save zero value
var once sync.Once // make sure code only run 1 time in entire program enventhough many goroutine callback
var configInstance *Config

func GetConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config")
		viper.SetConfigType("ymal")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		// Marshal = chuyển dữ liệu Go (struct, map, slice) thành JSON (dạng []byte).
		// Unmarshal = chuyển JSON thành dữ liệu Go.
		if err := viper.Unmarshal(&configInstance); err != nil {
			panic(err)
		}
	})

	return configInstance
}
