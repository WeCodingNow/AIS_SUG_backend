package main

import (
	"log"

	"github.com/WeCodingNow/AIS_SUG_backend/internal/server"
)

func main() {
	app := server.NewApp()
	app.Init()

	if err := app.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
