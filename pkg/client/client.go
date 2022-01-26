package client

import (
	"context"

	// "log"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	// "github.com/sonr-io/sonr/core/device"
	// "github.com/sonr-io/sonr/pkg/crypto"

	"github.com/sonr-io/sonr/pkg/p2p"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
	highwayv1 "go.buf.build/grpc/go/sonr-io/highway/v1"
	"google.golang.org/grpc"
)

var sonrDeviceKeyFile = "sonr-provision.priv"

// Client is a client for the Sonr network
type Client struct {
	cosmosclient.Client
	highwayv1.HighwayServiceClient
	Host    p2p.HostImpl
	ctx     context.Context
	Account cosmosaccount.Account
}

// NewClient creates a new client for the given name
func NewClient(ctx context.Context, addr string, sname string, passphrase string) (*Client, error) {
	// create an instance of cosmosclient
	cosmos, err := cosmosclient.New(ctx, cosmosclient.WithKeyringBackend(cosmosaccount.KeyringOS))
	if err != nil {
		return nil, err
	}

	// create an instance of highwayv1
	acc, _, err := cosmos.AccountRegistry.Create(sname)
	if err != nil {
		acc, err = cosmos.AccountRegistry.GetByName(sname)
		if err != nil {
			return nil, err
		}
	}

	// walletFolder, err := device.Support.CreateFolder(".wallet")
	// if err != nil {
	// 	return nil, err
	// }

	// create a new keyring
	// ks, m, err := crypto.GenerateKeyring(sname, cosmos.Context.Keyring, crypto.WithFolder(walletFolder))
	// if err != nil {
	// 	return nil, err
	// }
	// log.Printf("Generated keyring with mnemonic: %s \n", m)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	// // Create a new p2p host
	// h, err := p2p.NewHost(ctx, ks.CryptoPrivKey())
	// if err != nil {
	// 	return nil, err
	// }

	// Create a new highway client
	highway := highwayv1.NewHighwayServiceClient(conn)
	return &Client{
		Client:               cosmos,
		HighwayServiceClient: highway,
		ctx:                  ctx,
		Account:              acc,
		// Host:                 h,
	}, nil
}

// Keyring returns the keyring for the given name
func (c *Client) Keyring() keyring.Keyring {
	return c.Client.Context.Keyring
}
