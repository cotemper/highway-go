/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"context"
	"log"

	"github.com/sonr-io/highway-go/config"
	highway "github.com/sonr-io/highway-go/grpc"
)

func main() {

	//TODO cosmos setup
	// Creates a new client with an `.snr` name called 'test' and a wallet passphrase called 'bad-password'
	// and attempts to connect to the Sonr Blockchain node at the address "http://localhost:26657"
	//_, err := client.NewClient(context.Background(), "http://127.0.0.1:26657", "test", "bad-password")

	var err error = nil
	if err != nil {
		log.Fatal(err)
	}

	// This loads the information from the configuration file into a struct
	cnfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// This is starting the Highway based Node utilizing the configuration
	err = highway.Start(context.Background(), cnfg)
	if err != nil {
		log.Fatal(err)
	}

	// This is bad practice however this is being used to block the code
	// from completing execution
	//select {}
}
