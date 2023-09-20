package main

import (
	"awesomeProject/internal/api"
	"fmt"
	"log"
)

func main() {
	log.Println("Application start!")
	api.StartServer()
	fmt.Println("Application terminated!")
}
