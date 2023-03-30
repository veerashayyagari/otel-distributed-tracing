package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/veerashayyagari/go-otel/services/app/handlers"
)

var (
	build = "local"
	port  = "3000"
)

func main() {
	log.Printf("starting web app for build: %s  \n", build)
	defer log.Println("completed shutting down web app")

	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- http.ListenAndServe(fmt.Sprintf(":%s", port), handlers.AppRouter())
	}()

	select {
	case err := <-serverErrors:
		log.Printf("unhandled server error. %s \n", err)
	case v := <-shutdown:
		log.Printf("received shutdown signal %v. Shutting down web app.\n", v)
	}
}
