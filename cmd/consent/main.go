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
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	port := flag.Int("port", 8000, "The port to bind to")
	domain := flag.String("domain", "", "The domain used to serve the application via SSL")
	certs := flag.String("certs", "/var/www/.cache", "The directory to use for caching SSL certificates")
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
		if *domain != "" {
			m := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(*domain),
				Cache:      autocert.DirCache(*certs),
			}
			go http.ListenAndServe(":http", m.HTTPHandler(nil))
			if err := http.Serve(m.Listener(), srv.Handler); err != nil && err != http.ErrServerClosed {
				logger.Fatalf("cmd: error binding server to network: %s", err.Error())
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.Fatalf("cmd: error starting server: %s", err.Error())
			}
		}
	}()

	if *domain != "" {
		logger.Printf("Server is using AutoTLS and is now listening on ports 80 and 443")
	} else {
		logger.Printf("Server now listening on port %d", *port)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("cmd:  error shutting down server: %s", err.Error())
	}
}
