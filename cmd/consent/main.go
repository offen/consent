// Copyright 2022 - Offen Authors <hioffen@posteo.de>
// SPDX-License-Identifier: Apache-2.0

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
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffyaml"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	logger := log.New(os.Stderr, "", log.Ldate)

	fs := flag.NewFlagSet("consent", flag.ContinueOnError)
	var (
		port   = fs.Int("port", 8000, "The port to bind to (also via PORT)")
		domain = fs.String("domain", "", "The domain used to serve the application via SSL (also via DOMAIN)")
		certs  = fs.String("certs", "/var/www/.cache", "The directory to use for caching SSL certificates (also via CERTS)")

		copy      = fs.String("ui-copy", "", "The copy used for the default consent banner (also via UI_COPY)")
		buttonYes = fs.String("ui-button-yes", "", "The yes button used for the default consent banner (also via UI_BUTTON_YES)")
		buttonNo  = fs.String("ui-button-no", "", "The No button used for the default consent banner (also via UI_BUTTON_NO)")

		templatesDirectory = fs.String("templates-directory", "", "The location to look for custom templates (also via TEMPLATES_DIRECTORY)")
		stylesheet         = fs.String("stylesheet", "", "The location to look for a custom style sheet (also via STYLESHEET)")

		_ = fs.String("config", os.Getenv("CONFIG_FILE"), "The location of the config file in yaml format (optional)")
	)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix(), ff.WithConfigFileFlag("config"), ff.WithConfigFileParser(ffyaml.Parser)); err != nil {
		logger.Fatalf("cmd: error parsing configuration: %s", err.Error())
		os.Exit(1)
	}

	handler, err := consent.NewHandler(
		consent.WithLogger(logger),
		consent.WithCustomizedWording(*copy, *buttonYes, *buttonNo),
		consent.WithTemplatesDirectory(*templatesDirectory),
		consent.WithStylesheet(*stylesheet),
	)
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
