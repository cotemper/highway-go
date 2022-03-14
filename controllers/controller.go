package controller

import (
	"context"

	"github.com/sonr-io/highway-go/config"
	db "github.com/sonr-io/highway-go/database"
	"github.com/sonr-io/highway-go/models"
	rt "go.buf.build/grpc/go/sonr-io/sonr/registry"
)

// any other services required by http server will flow through here
type Controller struct {
	client      db.MongoClient
	privateKey  string
	highwayStub *models.HighwayStub
}

func New(mongoClient db.MongoClient, cnfg *config.SonrConfig, stub *models.HighwayStub) (*Controller, error) {
	return &Controller{
		client:      mongoClient,
		privateKey:  cnfg.SecretKey,
		highwayStub: stub,
	}, nil
}

func (ctrl *Controller) CheckName(ctx context.Context, name string) (bool, error) {
	result, err := ctrl.client.CheckName(name)
	if err != nil {
		return false, err
	}
	return result, nil
}

func (ctrl *Controller) InsertRecord(ctx context.Context, recordObj db.RecordNameObj) error {
	err := ctrl.client.StoreRecord(recordObj)
	if err != nil {
		return err
	}
	return nil
}

func (ctrl *Controller) RegisterName(ctx context.Context, ret *rt.MsgRegisterName) (*rt.MsgRegisterNameResponse, error) {

	return &rt.MsgRegisterNameResponse{}, nil
}
