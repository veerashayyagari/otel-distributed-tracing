package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/veerashayyagari/go-otel/services/user-api/handlers"
)

const name = "user-api"

var (
	build   = "local"
	version = "1.0"
	port    = "4000"
)

func main() {
	log.Printf("starting %s, version: %s, for build: %s \n", name, version, build)
	defer log.Println("shutdown complete for ", name)

	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}

	// register shutdown channel to be notified on termination
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- http.ListenAndServe(fmt.Sprintf(":%s", port), handlers.UserAPIRouter())
	}()

	select {
	case err := <-serverErrors:
		if err != nil {
			log.Fatalf("unhandled server error. %s \n", err)
		}
	case sig := <-shutdown:
		log.Printf("shutting down %s, received signal %v \n", name, sig)
	}
}
