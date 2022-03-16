package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sonr-io/webauthn.io/config"
	"github.com/sonr-io/webauthn.io/controller"
	db "github.com/sonr-io/webauthn.io/database"
	highway "github.com/sonr-io/webauthn.io/highway"
	"github.com/sonr-io/webauthn.io/logger"
	log "github.com/sonr-io/webauthn.io/logger"
	"github.com/sonr-io/webauthn.io/server"
)

func main() {
	highwayConfig, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	// authConfig, err := config.LoadConfig("config.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	authConfig := &config.Config{
		HostAddress:  highwayConfig.HighwayAddress,
		HostPort:     highwayConfig.HttpPort,
		DBName:       highwayConfig.SqlName,
		DBPath:       highwayConfig.SqlPath,
		RelyingParty: highwayConfig.RelyingParty,
	}

	err = log.Setup(authConfig)
	if err != nil {
		log.Fatal(err)
	}

	//get ctrl for highway
	stub := highway.Start(context.Background(), highwayConfig)
	DB, err := db.Connect(highwayConfig.MongoUri, highwayConfig.MongoCollectionName, highwayConfig.MongoDbName)
	if err != nil {
		logger.Errorf("database connection failed")
	}

	ctrl, err := controller.New(DB, highwayConfig, stub)
	if err != nil {
		log.Fatal(err)
	}

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
