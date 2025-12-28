package main

import (
	"log"
	"sihce_consulta_externa/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Error fatal: %v", err)
	}
}
