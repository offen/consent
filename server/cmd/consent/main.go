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

	consent "github.com/offen/consent/server"
)

func main() {
	port := flag.String("port", "8000", "The port to bind to")
	handler, err := consent.NewHandler()
	if err != nil {
		panic(fmt.Errorf("cmd: error creating handler: %w", err))
	}

	srv := &http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%s", *port),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("cmd: error starting server: %w", err))
		}
	}()

	log.Printf("Server now listening on port %s", *port)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		panic(fmt.Errorf("cmd: error shutting down server: %w", err))
	}
}
