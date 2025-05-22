package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/dewidyabagus/go-payout-workflow/sources/pkg/validator"
	"github.com/dewidyabagus/go-payout-workflow/sources/pkg/workflow"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	workflow := workflow.New(workflow.Config{
		BaseURL: "http://localhost:5000/conductor/api",
		Authorization: &workflow.BasicAuth{
			Username: "admin",
			Password: "qwerty",
		},
	})
	handler := handler{
		validator: validator.New(),
		service: &service{
			workflow: workflow.WorkflowExecutor(),
		},
	}
	r.Post("/ppob/prepaid-payment", handler.PrepaidPaymentHandler)
	r.Post("/webhooks/harsya/payment-notifications", handler.HarsyaPaymentNotificationHandler)

	signal, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	server := &http.Server{
		Addr:         ":7000",
		Handler:      r,
		ReadTimeout:  6 * time.Second,
		WriteTimeout: 90 * time.Second,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln("HTTP Server Error:", err.Error())
		}
	}()
	log.Println("HTTP service running on", server.Addr)

	<-signal.Done()

	stop()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("Failed while shutting down http server. Error:", err.Error())
	}
}
