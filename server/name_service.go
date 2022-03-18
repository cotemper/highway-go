package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	rt "go.buf.build/grpc/go/sonr-io/sonr/registry"
)

func (ws *Server) CheckName(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	vars := mux.Vars(req)
	//The trimmer
	name := vars["name"]
	if name[len(name)-4:] == ".snr" {
		name = name[:len(name)-4]
	}
	var err error

	// start := time.Now()
	// e := log.Info()
	// defer func(e *zerolog.Event, start time.Time) {
	// 	if err != nil {
	// 		e = log.Error().Stack().Err(err)
	// 	}
	// 	e.Str("handler", "CheckName").AnErr("context", ctx.Err()).Str("name", name).Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
	// }(e, start)

	nameAvailable, err := ws.Ctrl.CheckName(ctx, name)
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

//TODO clean up to match other calls
func (ws *Server) RegisterName(w http.ResponseWriter, req *http.Request) {
	//var body *rt.MsgRegisterName
	ctx := req.Context()
	var err error

	vars := mux.Vars(req)
	name := vars["name"]

	// start := time.Now()
	// e := log.Info()
	// defer func(e *zerolog.Event, start time.Time) {
	// 	if err != nil {
	// 		e = log.Error().Stack().Err(err)
	// 	}
	// 	e.Str("handler", "RegisterName").AnErr("context", ctx.Err()).Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
	// }(e, start)

	// body, err := ioutil.ReadAll(req.Body)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// }
	//log.Debug().Str("handler", "RegisterName").Bytes("request_body", body).Send()

	// var recObj *rt.MsgRegisterName
	// err = json.Unmarshal(body, &recObj)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// }

	//TODO checkname
	user := ws.Ctrl.FindUserByName(ctx, name)
	if user.Username == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	did := user.Did

	resp, err := ws.Ctrl.RegisterName(ctx, &rt.MsgRegisterName{NameToRegister: name}, did, nil)
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
