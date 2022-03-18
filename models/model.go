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
<<<<<<< HEAD

	// publickey.challenge.userID
	// user: {
	//         id: Uint8Array.from(
	//             "UZSL85T9AFC", c => c.charCodeAt(0)),
	//         name: "lee@webauthn.guide",
	//         displayName: "Lee",
	//     },
	//     pubKeyCredParams: [{alg: -7, type: "public-key"}],
	//     authenticatorSelection: {
	//         authenticatorAttachment: "cross-platform",
	//     },

=======
>>>>>>> 9ce0b9c53cf9b63af806fcb90bd2962ecbb0245d
}
