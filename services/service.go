package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"github.com/kataras/jwt"
	controller "github.com/sonr-io/highway-go/controllers"
)

func AddHandlers(r *mux.Router, ctrl *controller.Controller) {
	// hello handler
	r.HandleFunc("/health/", HealthHandler).Methods("GET").Schemes("http")
	// JWT handler
	r.HandleFunc("/generate/", GenerateJWT).Methods("GET").Schemes("http")
	// check name
	r.HandleFunc("/check-name/", CheckName).Methods("GET").Schemes("http")
	// record registered name
	r.HandleFunc("/check-name/", RecordName).Methods("POST").Schemes("http")
}

// Error Definitions //TODO this has been used twice, move it a layer back and call it instead
var (
	logger                 = golog.Default.Child("grpc/highway")
	ErrEmptyQueue          = errors.New("no items in Transfer Queue")
	ErrInvalidQuery        = errors.New("no SName or PeerID provided")
	ErrMissingParam        = errors.New("paramater is missing")
	ErrProtocolsNotSet     = errors.New("node Protocol has not been initialized")
	ErrMethodUnimplemented = errors.New("method is not implemented")
)

// Hello Handler
func HealthHandler(w http.ResponseWriter, req *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
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

func CheckName(w http.ResponseWriter, req *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func RecordName(w http.ResponseWriter, req *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}
