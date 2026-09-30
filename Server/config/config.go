package config

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	Server            *Server
	Db                *Db
	JWTPrivateKeyPath string
	JWTPublicKeyPath  string
	CookieSecure      bool
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
		workingDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		configPath := findEnvFile(workingDir)
		viper.SetConfigFile(configPath)
		viper.SetConfigType("env")
		viper.AutomaticEnv()

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}

		configInstance = &Config{
			JWTPrivateKeyPath: resolveConfigPath(configPath, viper.GetString("JWT_PRIVATE_KEY_PATH")),
			JWTPublicKeyPath:  resolveConfigPath(configPath, viper.GetString("JWT_PUBLIC_KEY_PATH")),
			CookieSecure:      viper.GetBool("COOKIE_SECURE"),
			Server: &Server{
				Port: viper.GetInt("PORT"),
			},
			Db: &Db{
				Host:     viper.GetString("DB_HOST"),
				Port:     viper.GetInt("DB_PORT"),
				User:     viper.GetString("DB_USER"),
				Password: viper.GetString("DB_PASSWORD"),
				DBName:   viper.GetString("DB_NAME"),
				SSLMode:  viper.GetString("DB_SSL_MODE"),
				TimeZone: viper.GetString("DB_TIME_ZONE"),
			},
		}
	})

	return configInstance
}

func resolveConfigPath(configPath, configuredPath string) string {
	if filepath.IsAbs(configuredPath) {
		return configuredPath
	}
	return filepath.Join(filepath.Dir(configPath), configuredPath)
}

func findEnvFile(startDir string) string {
	for directory := startDir; ; directory = filepath.Dir(directory) {
		configPath := filepath.Join(directory, ".env")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		parentDirectory := filepath.Dir(directory)
		if parentDirectory == directory {
			panic(".env file not found")
		}
	}
}
