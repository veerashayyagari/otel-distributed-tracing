package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/veerashayyagari/go-otel/services/user-api/handlers"
)

var (
	build = "local"
)

func main() {
	log.Printf("starting user api for build: %s \n", build)
	defer log.Println("shutdown complete for user api")

	// register shutdown channel to be notified on termination
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- http.ListenAndServe(":3000", handlers.UserAPIRouter())
	}()

	select {
	case err := <-serverErrors:
		log.Printf("received an unhandled server error: %s \n", err)
	case sig := <-shutdown:
		log.Printf("shutting down user-api, received signal %v \n", sig)
	}
}
