package controller

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/sonr-io/highway-go/config"
	db "github.com/sonr-io/highway-go/database"
	"github.com/sonr-io/highway-go/models"
	"github.com/sonr-io/sonr/x/registry/types"
	rt "go.buf.build/grpc/go/sonr-io/sonr/registry"
)

// any other services required by http server will flow through here
type Controller struct {
	client      db.MongoClient
	privateKey  string
	devAccount  string
	highwayStub *models.HighwayStub
}

func New(mongoClient db.MongoClient, cnfg *config.SonrConfig, stub *models.HighwayStub) (*Controller, error) {
	return &Controller{
		client:      mongoClient,
		privateKey:  cnfg.SecretKey,
		devAccount:  cnfg.DevAccount,
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

func (ctrl *Controller) RegisterName(ctx context.Context, req *rt.MsgRegisterName) (*rt.MsgRegisterNameResponse, error) {
	// account `alice` was initialized during `starport chain serve`
	//accountName := req.Creator
	accountName := ctrl.devAccount

	// get account from the keyring by account name and return a bech32 address
	address, err := ctrl.highwayStub.Cosmos.Address(accountName)
	if err != nil {
		fmt.Errorf("Error in cosmos.address: ")
		return &rt.MsgRegisterNameResponse{}, err
	}

	// define a message to create a post
	msg := &types.MsgRegisterName{
		Creator: address.String(),
		//DeviceId:       req.DeviceId,
		NameToRegister: req.NameToRegister,
		//Jwt:            req.PublicKey, //TODO implement new jwt system
	}

	fmt.Println(msg.NameToRegister)

	// broadcast a transaction from account accountName with the message to create a post
	//store response in txResp
	txResp, err := ctrl.highwayStub.Cosmos.BroadcastTx(req.Creator, msg)
	if err != nil {
		fmt.Errorf("Error in broadcastTx")
		return &rt.MsgRegisterNameResponse{}, err
	}

	//TODO fix this logic
	success := false
	if !txResp.Empty() {
		success = true
	}

	bs, err := hex.DecodeString(txResp.Data)
	if err != nil {
		fmt.Errorf("Error in hex.DecodeString")
		return &rt.MsgRegisterNameResponse{}, err
	}

	// Unmarshalling of a json did document:
	// parsedDIDDoc := did.Document{}
	// err = json.Unmarshal([]byte(bs), &parsedDIDDoc)
	// if err != nil {
	// 	fmt.Errorf("Error in json.Unmarshal ")
	// 	return &rt.MsgRegisterNameResponse{}, err
	// }
	//TODO unmarshal is not working as intended

	//did := "did:sonr:" + JWT
	response := rt.MsgRegisterNameResponse{}
	response.IsSuccess = success
	response.DidDocumentJson = string(bs)
	response.DidUrl = txResp.TxHash
	// response.DidDocument = &rt.DidDocument{
	// 	AlsoKnownAs: aliases,
	// }

	//return &response, nil
	return &response, nil
}
