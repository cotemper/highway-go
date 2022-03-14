package models

import (
	"context"
	"net"
	"net/http"

	"github.com/sonr-io/sonr/pkg/p2p"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
	hw "go.buf.build/grpc/go/sonr-io/highway/v1"
	"go.buf.build/grpc/go/sonr-io/sonr/channel"

	"google.golang.org/grpc"
)

// HighwayStub is the RPC Service for the Custodian Node.
type HighwayStub struct {
	hw.HighwayServer
	Host   p2p.HostImpl
	Cosmos cosmosclient.Client

	// Properties
	Ctx      context.Context
	Grpc     *grpc.Server
	Http     *http.Server
	Listener net.Listener

	// Configuration

	// List of Entries
	Channels map[string]channel.Channel
}

//TODO this needs work, fix this model, wtf are either of these fields???
type Jwt struct {
	Snr        string `json:"snr"`
	EthAddress string `json: "ethAddress"`
}

type User struct {
	Did string
	Jwt Jwt
}
