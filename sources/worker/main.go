package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/dewidyabagus/go-payout-workflow/sources/pkg/workflow"
)

func main() {
	signal, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	workflow := workflow.New(workflow.Config{
		BaseURL: "http://localhost:5000/conductor/api",
		Authorization: &workflow.BasicAuth{
			Username: "admin",
			Password: "qwerty",
		},
	})
	handler := &handler{
		service: &service{},
	}

	workers := workflow.Workers()
	defer workers.Close()

	go func() {
		if err := workers.RunWorkers(workerDefinitions(handler)); err != nil {
			log.Fatalln("Run workers:", err.Error())
		}
	}()

	<-signal.Done()

	stop()
	log.Println("Graceful shutdown started ...")
}
