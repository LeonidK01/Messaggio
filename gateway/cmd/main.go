package main

import (
	"log"

	"github.com/LeonidK01/Messaggio/internal/app"
)

func main() {
	if err := app.Start(); err != nil {
		log.Fatal("failed start application")
	}
}
