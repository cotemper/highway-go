package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kataras/golog"
	"github.com/kataras/jwt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	controller "github.com/sonr-io/highway-go/controllers"
	db "github.com/sonr-io/highway-go/database"
	rt "go.buf.build/grpc/go/sonr-io/sonr/registry"
)

// TODO expand with some kind of middleware later
func AddHandlers(r *mux.Router, ctrl *controller.Controller) {
	// hello handler
	r.HandleFunc("/health", HealthHandler(ctrl)).Methods("GET").Schemes("http")
	// JWT handler
	r.HandleFunc("/generate", GenerateJWT(ctrl)).Methods("GET").Schemes("http")
	// check name
	r.HandleFunc("/check/name/{name}", CheckName(ctrl)).Methods("GET").Schemes("http")
	// record registered name
	r.HandleFunc("/record/name", RecordName(ctrl)).Methods("POST").Schemes("http")
	// Register a name
	r.HandleFunc("/register/name", RegisterName(ctrl)).Methods("POST").Schemes("http")
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

// GenerateJWT generates a JWT for the given SName and PeerID.
func GenerateJWT(ctrl *controller.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		keys, ok := req.URL.Query()["token"]
		if !ok || len(keys[0]) < 1 {
			logger.Warn("Url Param 'key' is missing")
			return
		}

		token := keys[0]

		fmt.Println(token)

		verifiedToken, err := jwt.Verify(jwt.HS256, []byte(""), []byte(token))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result := Jwt{}
		err = verifiedToken.Claims(&result)
		if err != nil {
			logger.Fatalf("JWT Error", err)
		}

		//format response
		js, err := json.Marshal(result)
		if err != nil {
			logger.Fatalf("Error happened in JSON marshal. Err: %s", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

type Response struct {
	Available bool
}

func CheckName(ctrl *controller.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		vars := mux.Vars(req)
		name := vars["name"]
		var err error

		//TODO add error checking for bad values

		start := time.Now()
		e := log.Info()
		defer func(e *zerolog.Event, start time.Time) {
			if err != nil {
				e = log.Error().Stack().Err(err)
			}
			e.Str("handler", "CheckName").AnErr("context", ctx.Err()).Str("name", name).Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
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
		ctx := req.Context()
		var err error

		start := time.Now()
		e := log.Info()
		defer func(e *zerolog.Event, start time.Time) {
			if err != nil {
				e = log.Error().Stack().Err(err)
			}
			e.Str("handler", "RecordName").AnErr("context", ctx.Err()).Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
		}(e, start)

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		log.Debug().Str("handler", "recordName").Bytes("request_body", body).Send()

		var recObj db.RecordNameObj
		err = json.Unmarshal(body, &recObj)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		available, err := ctrl.CheckName(ctx, recObj.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if !available {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = ctrl.InsertRecord(ctx, recObj)

		if err != nil {
			w.WriteHeader(http.StatusExpectationFailed)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func RegisterName(ctrl *controller.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		//var body *rt.MsgRegisterName
		ctx := req.Context()
		var err error

		start := time.Now()
		e := log.Info()
		defer func(e *zerolog.Event, start time.Time) {
			if err != nil {
				e = log.Error().Stack().Err(err)
			}
			e.Str("handler", "RegisterName").AnErr("context", ctx.Err()).Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
		}(e, start)

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		log.Debug().Str("handler", "RegisterName").Bytes("request_body", body).Send()

		var recObj *rt.MsgRegisterName
		err = json.Unmarshal(body, &recObj)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		resp, err := ctrl.RegisterName(ctx, recObj)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		//format response
		js, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
