package main

import (
	"log"

	"github.com/WeCodingNow/AIS_SUG_backend/server"
)

func main() {
	app := server.NewApp()

	if err := app.Run("8080"); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
