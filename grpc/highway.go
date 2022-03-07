package highway

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/kataras/golog"
	"github.com/phayes/freeport"
	"github.com/sonr-io/highway-go/config"
	"github.com/sonr-io/sonr/pkg/p2p"

	channel "github.com/sonr-io/sonr/x/channel/service"
	hw "go.buf.build/grpc/go/sonr-io/highway/v1"

	"github.com/tendermint/starport/starport/pkg/cosmosclient"
	"google.golang.org/grpc"
)

// Error Definitions
var (
	logger                 = golog.Default.Child("grpc/highway")
	ErrEmptyQueue          = errors.New("No items in Transfer Queue.")
	ErrInvalidQuery        = errors.New("No SName or PeerID provided.")
	ErrMissingParam        = errors.New("Paramater is missing.")
	ErrProtocolsNotSet     = errors.New("Node Protocol has not been initialized.")
	ErrMethodUnimplemented = errors.New("Method is not implemented.")
)

// HighwayStub is the RPC Service for the Custodian Node.
type HighwayStub struct {
	hw.HighwayServiceServer
	Host   p2p.HostImpl
	cosmos cosmosclient.Client

	// Properties
	ctx      context.Context
	grpc     *grpc.Server
	http     *http.Server
	listener net.Listener

	// Configuration

	// List of Entries
	channels map[string]channel.Channel
}

func Start(ctx context.Context, cnfg *config.SonrConfig) (*HighwayStub, error) {
	// Create the main listener.
	l, err := net.Listen(verifyAddress(cnfg))
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	print("\n\n\n\n\n\n\n\n\n\n\n")
	fmt.Print(l.Addr().Network())
	print("\n")
	fmt.Print(l.Addr().String())
	print("\n\n\n\n\n\n\n\n\n\n\n")

	// Create a cmux.
	stub, err := NewHighwayRPC(ctx, l, nil)
	if err != nil {
		return nil, err
	}
	return stub, nil
}

// NewHighway creates a new Highway service stub for the node.
func NewHighwayRPC(ctx context.Context, l net.Listener, h p2p.HostImpl) (*HighwayStub, error) {
	// // create an instance of cosmosclient
	// cosmos, err := cosmosclient.New(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// Create the RPC Service
	stub := &HighwayStub{
		Host: h,
		ctx:  ctx,
		grpc: grpc.NewServer(),
		//	cosmos:   cosmos,
		listener: l,
	}

	hw.RegisterHighwayServiceServer(stub.grpc, stub)
	go stub.Serve(ctx, l)
	return stub, nil
}

// Serve serves the RPC Service on the given port.
func (s *HighwayStub) Serve(ctx context.Context, listener net.Listener) {
	logger.Infof("Starting RPC Service on %s", listener.Addr().String())
	for {
		// Stop Serving if context is done
		select {
		case <-ctx.Done():
			// s.node.Close()
			return
		}
	}
}

// verifyAddress verifies the address is valid.
func verifyAddress(cnfg *config.SonrConfig) (string, string) {
	// Define Variables
	var network string
	var address string
	var port int
	var err error

	// Set Network
	if cnfg.HighwayNetwork != "tcp" {
		network = "tcp"
	}
	logger.Debugf("Network: %s", network)

	// Get free port if set port is 0 or 69420
	if cnfg.HighwayPort == 0 || cnfg.HighwayPort == 69420 {
		port, err = freeport.GetFreePort()
		if err != nil {
			panic(err)
		}
	}
	logger.Debugf("Using port %d", port)

	// Set Address
	if !strings.Contains(cnfg.HighwayAddress, ":") {
		address = fmt.Sprintf("%s:%d", cnfg.HighwayAddress, port)
	}
	logger.Debugf("Address: %s", address)

	return network, address
}
