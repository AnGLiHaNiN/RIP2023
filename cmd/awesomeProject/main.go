package main

import (
	"github.com/joho/godotenv"
	_ "gorm.io/driver/postgres"
	"log"
	"pharmaBlend/internal/app"
)

func main() {
	godotenv.Load()
	log.Println("Application start!")
	var config = &app.Config{
		LocalHost: "127.0.0.1",
		Port:      "8080",
	}

	// Pass both the config and repository to app.New
	a, err := app.New(config)
	if err != nil {
		log.Fatal("Failed to initialize the application:", err)
	}

	a.StartServer()
	log.Println("Application terminated!")
}
