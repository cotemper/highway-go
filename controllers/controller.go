package controller

import (
	db "github.com/sonr-io/highway-go/database"
)

// any other services required by http server will flow through here
type Controller struct {
	client db.MongoClient
}

func New(client db.MongoClient) *Controller {
	return &Controller{
		client: client,
	}
}
