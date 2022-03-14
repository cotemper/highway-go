package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"github.com/kataras/jwt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	controller "github.com/sonr-io/highway-go/controllers"
)

// TODO expand with some kind of middleware later
func AddHandlers(r *mux.Router, ctrl *controller.Controller) {
	// hello handler
	r.HandleFunc("/health", HealthHandler(ctrl)).Methods("GET").Schemes("http")
	// JWT handler
	r.HandleFunc("/generate", GenerateJWT(ctrl)).Methods("GET").Schemes("http")
	// check name
	r.HandleFunc("/check-name/{name}", CheckName(ctrl)).Methods("GET").Schemes("http")
	// record registered name
	r.HandleFunc("/record-name", RecordName(ctrl)).Methods("POST").Schemes("http")
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
func HealthHandler(ctrl *controller.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// A very simple health check.
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
	}
}

type Jwt struct {
	Snr        string `json:"snr"`
	EthAddress string `json: "ethAddress"`
}

// Keep it secret.
var sharedKey = os.Getenv("FAKEPASSWORD")

// GenerateJWT generates a JWT for the given SName and PeerID.
func GenerateJWT(ctrl *controller.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
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
}

type Response struct {
	Available bool
}

func CheckName(ctrl *controller.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		vars := mux.Vars(req)
		name := vars["id"]
		var err error

		//TODO add error checking for bad values

		//var err error
		start := time.Now()
		e := log.Info()
		defer func(e *zerolog.Event, start time.Time) {
			if err != nil {
				e = log.Error().Stack().Err(err)
			}
			e.Str("handler", "getLineItem").AnErr("context", ctx.Err()).Str("name", name).Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
		}(e, start)

		nameAvailable, err := ctrl.CheckName(ctx, name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//format response
		responseObj := Response{Available: nameAvailable}
		js, err := json.Marshal(responseObj)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}

func RecordName(ctrl *controller.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// A very simple health check.
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
	}
}
