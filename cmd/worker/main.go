package main

import (
	"log"
	"shortener/internal/machinery"
)

func main() {
	server := machinery.New()
	worker := server.NewWorker("email_worker", 10)
	err := worker.Launch()
	if err != nil {
		log.Fatal(err)
		return
	}
}
