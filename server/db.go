package server

import (
	"database/sql"
	"fmt"

	"github.com/WeCodingNow/AIS_SUG_backend/config"
	"github.com/spf13/viper"
)

func makePostgresString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("postgres.host"), viper.GetInt("postgres.port"),
		viper.GetString("postgres.user"), viper.GetString("postgres.password"),
		viper.GetString("postgres.dbname"),
	)
}

func initDB() *sql.DB {
	if err := config.Init(); err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", makePostgresString())

	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
