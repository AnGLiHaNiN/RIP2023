package main

import (
	"R_I_P_labs/internal/pkg/app"
	"log"
)

// TODO: change
// @title Electronic notifications
// @version 1.0

// @host 127.0.0.1:8000
// @schemes http
// @BasePath /

func main() {
	app, err := app.New()
	if err != nil {
		log.Println("app can not be created", err)
		return
	}
	app.Run()
}
