package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/veerashayyagari/go-otel/services/product-api/handlers"
)

var (
	build = "local"
	port  = "6000"
)

func main() {
	log.Printf("starting product api for build: %s \n", build)
	defer log.Println("shutdown complete for product api")

	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}

	// register shutdown channel to be notified on termination
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- http.ListenAndServe(fmt.Sprintf(":%s", port), handlers.ProductAPIRouter())
	}()

	select {
	case err := <-serverErrors:
		log.Printf("received an unhandled server error: %s \n", err)
	case sig := <-shutdown:
		log.Printf("shutting down product-api, received signal %v \n", sig)
	}
}
