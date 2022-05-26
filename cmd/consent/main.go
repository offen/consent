package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	consent "github.com/offen/consent"
)

func main() {
	port := flag.Int("port", 8000, "The port to bind to")
	flag.Parse()

	logger := log.New(os.Stderr, "", log.Ldate)

	handler, err := consent.NewHandler(consent.WithLogger(logger))
	if err != nil {
		logger.Fatalf("cmd: error creating handler: %s", err.Error())
	}

	srv := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%d", *port),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("cmd: error starting server: %s", err.Error())
		}
	}()

	logger.Printf("Server now listening on port %d", *port)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("cmd:  error shutting down server: %s", err.Error())
	}
}
