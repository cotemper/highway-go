package models

import (
	"context"
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
	Ctx  context.Context
	Grpc *grpc.Server
	Http *http.Server

	// Configuration

	// List of Entries
	Channels map[string]channel.Channel
}

//get
// no clear answer

//give
//did

//TODO this needs work, remove soon
type Jwt struct {
	Snr        string `json:"snr"`
	EthAddress string `json: "ethAddress"`
}
