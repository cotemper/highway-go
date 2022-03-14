package controller

import (
	"context"

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

func (ctrl *Controller) CheckName(ctx context.Context, name string) (bool, error) {

	result, err := ctrl.client.CheckName(name)
	if err != nil {
		return false, err
	}
	return result, nil
}
