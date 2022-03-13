package highway

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"github.com/kataras/jwt"
	"github.com/sonr-io/highway-go/config"
	"github.com/sonr-io/highway-go/reflection"
	"github.com/sonr-io/sonr/pkg/p2p"
	"google.golang.org/grpc/credentials"

	channel "github.com/sonr-io/sonr/x/channel/service"
	hw "go.buf.build/grpc/go/sonr-io/highway/v1"

	"github.com/sonr-io/highway-go/pkg/client"
	"github.com/tendermint/starport/starport/pkg/cosmosclient"
	"google.golang.org/grpc"
)

const (
	// PEM_CERT_FILE is the path to the certificate file.
	PEM_CERT_FILE = "cert.pem"

	// PEM_KEY_FILE is the file containing the private key.
	PEM_KEY_FILE = "key.pem"
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
	hw.HighwayServer
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

// Hello Handler
func HelloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

type Jwt struct {
	Snr        string `json:"snr"`
	EthAddress string `json: "ethAddress"`
}

// Keep it secret.
var sharedKey = os.Getenv("FAKEPASSWORD")

// GenerateJWT generates a JWT for the given SName and PeerID.
func GenerateJWT(w http.ResponseWriter, req *http.Request) {
	keys, ok := req.URL.Query()["token"]
	if !ok || len(keys[0]) < 1 {
		logger.Warn("Url Param 'key' is missing")
		return
	}

	tokenString := keys[0]
	verifiedToken, err := jwt.Verify(jwt.HS256, sharedKey, []byte(tokenString))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := Jwt{}
	err = verifiedToken.Claims(&result)
	if err != nil {
		logger.Fatalf("JWT Error", err)
	}

	resp := make(map[string]string)
	resp["message"] = "Status Created"
	jsonResp, err := json.Marshal(result)
	if err != nil {
		logger.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}

// Start starts the RPC Service.
func Start(ctx context.Context, cnfg *config.SonrConfig) error {
	r := mux.NewRouter()
	// hello handler
	r.HandleFunc("/", HelloHandler)
	// file handler
	r.HandleFunc("/generate/", GenerateJWT).Methods("GET").Schemes("http")

	//get http port
	httpAddr := cnfg.HttpPort

	// Check if files exists
	if fileExists(PEM_CERT_FILE) && fileExists(PEM_KEY_FILE) {
		logger.Infof("Using TLS")
		go http.ListenAndServeTLS(":"+httpAddr, PEM_CERT_FILE, PEM_KEY_FILE, r)
	} else {
		logger.Warn("Using insecure HTTP")
		go http.ListenAndServe(":"+httpAddr, r)
	}

	// Create the main listener.
	l, err := net.Listen(verifyAddress(cnfg))
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	logger.Infof("Network: " + l.Addr().Network())
	logger.Infof("Address: " + l.Addr().String())

	// create an instance of cosmosclient
	cosmos, err := client.NewClient(context.Background(), l.Addr().String(), "test", "unimplemented-password")
	if err != nil {
		logger.Fatal("your cosmos is bad") //TODO error better when you're done debugging
		return err
	}

	// Get TLS config if TLS is enabled
	var stub *HighwayStub
	credentials, err := loadTLSCredentials()
	if err != nil {
		logger.Infof("Error loading TLS credentials: ", err)

		// If TLS is not enabled, create a new listener.
		// Create the RPC Service
		stub = &HighwayStub{
			Host:     nil,
			ctx:      ctx,
			grpc:     grpc.NewServer(),
			cosmos:   cosmos.Client,
			listener: l,
		}
		hw.RegisterHighwayServer(stub.grpc, stub)
		reflection.RegisterReflection(stub.grpc)
		logger.Infof("Starting RPC Service on %s", l.Addr().String())
		return stub.grpc.Serve(l)
	}

	// Create the RPC Service
	stub = &HighwayStub{
		Host:     nil,
		ctx:      ctx,
		grpc:     grpc.NewServer(grpc.Creds(credentials)),
		cosmos:   cosmos.Client,
		listener: l,
	}
	hw.RegisterHighwayServer(stub.grpc, stub)
	reflection.RegisterReflection(stub.grpc)
	logger.Infof("Starting RPC Service on %s", l.Addr().String())
	return stub.grpc.Serve(l)
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
	logger.Infof("Network: %s", network)

	//set port
	port, err = strconv.Atoi(cnfg.GrpcPort)
	if err != nil {
		return "", err.Error()
	}
	logger.Infof("Using port %d", port)

	// Set Address
	if !strings.Contains(cnfg.HighwayAddress, ":") {
		address = fmt.Sprintf("%s:%d", cnfg.HighwayAddress, port)
	}
	logger.Infof("Address: %s", address)

	return network, address
}

// Helper function to see if file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(PEM_CERT_FILE, PEM_KEY_FILE)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}
