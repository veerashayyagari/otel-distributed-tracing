package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/veerashayyagari/go-otel/services/app/router"
	"github.com/veerashayyagari/go-otel/tracer"
)

const name = "webapp"

var (
	build   = "local"
	version = "1.0"
	port    = "3000"
)

func main() {
	log.Printf("starting %s, version: %s, for build: %s  \n", name, version, build)
	defer log.Println("completed shutting down ", name)

	if p, ok := os.LookupEnv("PORT"); ok {
		port = p
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	go func() {
		cfg := &tracer.TraceConfig{
			ServiceName:    name,
			ServiceVersion: version,
			Environment:    build,
			ExportURI:      os.Getenv("ZIPKIN_API_URI"),
		}

		tr, err := tracer.NewTraceProvider(cfg)
		if err != nil {
			log.Println("failed to setup tracer.", err)
		}

		serverErrors <- http.ListenAndServe(fmt.Sprintf(":%s", port), router.New(tr))
	}()

	select {
	case err := <-serverErrors:
		if err != nil {
			log.Fatalf("unhandled server error. %s \n", err)
		}

	case v := <-shutdown:
		log.Printf("received shutdown signal %v. Shutting down %s.\n", v, name)
	}
}
