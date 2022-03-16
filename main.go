package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sonr-io/webauthn.io/config"
	highway "github.com/sonr-io/webauthn.io/highway"
	log "github.com/sonr-io/webauthn.io/logger"
	"github.com/sonr-io/webauthn.io/models"
	"github.com/sonr-io/webauthn.io/server"
)

func main() {
	highwayConfig, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	authConfig, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	err = models.Setup(authConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = log.Setup(authConfig)
	if err != nil {
		log.Fatal(err)
	}

	//get ctrl for highway
	ctrl, err := highway.Start(context.Background(), highwayConfig)

	// start server
	server, err := server.NewServer(ctrl, authConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	go server.Start()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGINT)

	<-c
	log.Info("Shutting down...")
	server.Shutdown()
}
