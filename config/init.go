package config

import (
	"github.com/spf13/viper"
)

func initServer() {
	viper.SetDefault("server", map[string]interface{}{
		"host": "localhost",
		"port": 8080,
	})
}

func initPostgres() {
	viper.SetDefault("postgres", map[string]interface{}{
		"host":     "localhost",
		"port":     5432,
		"user":     "postgres",
		"password": "postgres",
		"dbname":   "postgres",
	})
}

func initAuth() {
	viper.SetDefault("auth", map[string]interface{}{
		"hash_salt":   "12345678",
		"signing_key": "87654321",
		"token_ttl":   86400,
	})
}

func Init() error {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	initPostgres()
	initAuth()
	initServer()

	return viper.ReadInConfig()
}
