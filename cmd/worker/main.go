package main

import (
	"log"
	"shortener/internal"
	"shortener/internal/machinery"
)

func main() {
	config := internal.MustLoad()
	server := machinery.New(config.Mail.Email, config.Mail.Password)
	worker := server.NewWorker("email_worker", 10)
	err := worker.Launch()
	if err != nil {
		log.Fatal(err)
		return
	}
}
