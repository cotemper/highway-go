package main

import (
	"context"
	"log"

	highway "github.com/sonr-io/highway-go/lib"
	"github.com/sonr-io/highway-go/pkg/client"
	"github.com/sonr-io/sonr/config"
)

func main() {
	_, err := client.NewClient(context.Background(), "http://localhost:26657", "test", "bad-password")
	if err != nil {
		log.Fatal(err)
	}
	cnfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	_, err = highway.Start(context.Background(), cnfg)
	if err != nil {
		log.Fatal(err)
	}
}
