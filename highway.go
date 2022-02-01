package main

import (
	"context"
	"log"

	"github.com/sonr-io/highway-go/pkg/client"
)

func main() {
	_, err := client.NewClient(context.Background(), "http://localhost:26657", "test", "bad-password")
	if err != nil {
		log.Fatal(err)
	}

}
