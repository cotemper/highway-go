package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	rt "go.buf.build/grpc/go/sonr-io/sonr/registry"
)

func (ws *Server) CheckName(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var takenNames = []string{"api", "tx", "app", "arianagrande", "azamsharp", "barrybonds", "barrysanders", "billgates", "britneyspears", "cdixon", "cristiano", "drake", "elon", "eminem", "flotus", "iamsrk", "imap", "index", "jack", "jbbernstein", "jeffbezos", "jimmyfallon", "joynerlucas", "jtimberlake", "justinbieber", "katyperry", "kimkardashian", "kingjames", "ladygaga", "larrypage", "launchhouse", "logic", "mail", "main", "markzuckerburg", "meekmill", "naval", "neymarjr", "oprah", "patrickbetdavid", "pop", "potus", "prad", "rihanna", "root", "satyanadella", "sc", "selenagomez", "sergeibrin", "shakira", "shl", "smartrick", "srbachchan", "stephencurry", "sundarpichai", "taylorswift", "tombrady", "vitalik", "michael", "prad2", "papa", "ikj", "ian", "shadowysupercoder", "ianperez", "perez", "0x0", "zac", "smartrick", "holwerda", "zholwerda", "NFT", "classof.o7", "goat", "nsfw", "nick", "ntindle", "nicktindle", "cloud", "devops", "engineer", "ntt", "grace", "get", "gtindle", "0xDEADBEEF", "static", "d0x", "null", "exposure", "zach", "joshLong145", "beanPole", "undefined", "Peyton", "gopher", "cosmic", "lauren", "sonr", "prad", "letsgobrandon", "snr", "erin", "jamey", "monica", "Space", "timmy", "creaton", "Warriors", "BestButt", "Mfers", "Beast", "mary", "david", "RX", "NT", "0X", "OK", "NO", "SN", "GB", "GT", "IP", "AH", "PT", "JL", "AF", "0F", "0p", "00", "C0", "80"}
	vars := mux.Vars(req)
	name := vars["name"]
	//The trimmer
	if name[len(name)-4:] == ".snr" {
		name = name[:len(name)-4]
	}
	var err error

	// if reserved
	for _, x := range takenNames {
		if x == name {
			//TODO return error
			return
		}
	}

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
