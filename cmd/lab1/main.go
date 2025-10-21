package main

import (
	"log"

	"lab1/internal/api"
)

func main() {
	log.Println("Application start up")
	api.StartServer()
	log.Println("Application down")
}

