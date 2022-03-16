package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/mux"
	log "github.com/sonr-io/webauthn.io/logger"
	"github.com/sonr-io/webauthn.io/models"
	rt "go.buf.build/grpc/go/sonr-io/sonr/registry"
)

// RequestNewCredential begins a Credential Registration Request, returning a
// PublicKeyCredentialCreationOptions object
func (ws *Server) RequestNewCredential(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vars := mux.Vars(r)
	username := vars["name"]

	// Most times relying parties will choose these.
	attType := r.FormValue("attType")
	authType := r.FormValue("authType")

	// Advanced settings
	userVer := r.FormValue("userVerification")
	resKey := r.FormValue("residentKeyRequirement")
	testExtension := r.FormValue("txAuthExtension")

	var residentKeyRequirement *bool
	if strings.EqualFold(resKey, "true") {
		residentKeyRequirement = protocol.ResidentKeyRequired()
	} else {
		residentKeyRequirement = protocol.ResidentKeyUnrequired()
	}

	testEx := protocol.AuthenticationExtensions(map[string]interface{}{"txAuthSimple": testExtension})

	//SQL lite check
	// user, err := models.GetUserByUsername(username)
	// if err != nil {
	// 	user = models.User{
	// 		DisplayName: strings.Split(username, "@")[0],
	// 		Username:    username,
	// 	}
	// 	err = models.PutUser(&user)
	// 	if err != nil {
	// 		jsonResponse(w, "Error creating new user", http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	//secondary mongo check
	user := ws.Ctrl.FindUserByName(ctx, username)

	// user doesn't exist, create new user
	if user.DisplayName == "" {
		available, _ := ws.Ctrl.CheckName(ctx, username)
		if !available {
			jsonResponse(w, fmt.Errorf("username is not availabel to use"), http.StatusAlreadyReported)
			return
		}
		var names []string
		names = append(names, username)
		did := "did:sonr:temp" + username
		user.DisplayName = username
		user.Names = names
		user.Did = did
		ws.Ctrl.NewUser(ctx, *user)
	}

	credentialOptions, sessionData, err := ws.webauthn.BeginRegistration(user,
		webauthn.WithAuthenticatorSelection(
			protocol.AuthenticatorSelection{
				AuthenticatorAttachment: protocol.AuthenticatorAttachment(authType),
				RequireResidentKey:      residentKeyRequirement,
				UserVerification:        protocol.UserVerificationRequirement(userVer),
			}),
		webauthn.WithConveyancePreference(protocol.ConveyancePreference(attType)),
		webauthn.WithExtensions(testEx),
	)
	if err != nil {
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Save the session data as marshaled JSON
	err = ws.store.SaveWebauthnSession("registration", sessionData, r, w)
	if err != nil {
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the PublicKeyCreationOptions back to the browser
	jsonResponse(w, credentialOptions, http.StatusOK)
	return
}

// MakeNewCredential attempts to make a new credential given an authenticator's response
func (ws *Server) MakeNewCredential(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Load the session data
	sessionData, err := ws.store.GetWebauthnSession("registration", r)
	if err != nil {
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Get the user associated with the credential
	user, err := ws.Ctrl.GetUser(models.BytesToID(sessionData.UserID))
	if err != nil {
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//get mongo user
	mgoUser := ws.Ctrl.FindUserByName(ctx, user.DisplayName)
	if mgoUser.DisplayName == "" {
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify that the challenge succeeded
	cred, err := ws.webauthn.FinishRegistration(user, sessionData, r)
	if err != nil {
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// If needed, you can perform additional checks here to ensure the
	// authenticator and generated credential conform to your requirements.

	// Finally, save the credential and authenticator to the
	// database
	authenticator, err := ws.Ctrl.CreateAuthenticator(cred.Authenticator)
	if err != nil {
		log.Errorf("error creating authenticator: %v", err)
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// For our use case, we're encoding the raw credential ID as URL-safe
	// base64 since we anticipate rendering it in templates. If you choose to
	// do this, make sure to decode the credential ID before passing it back to
	// the webauthn library.
	credentialID := base64.URLEncoding.EncodeToString(cred.ID)
	c := &models.Credential{
		Authenticator:   authenticator,
		AuthenticatorID: authenticator.ID,
		UserID:          user.ID,
		PublicKey:       cred.PublicKey,
		CredentialID:    credentialID,
	}
	err = ws.Ctrl.CreateCredential(c)
	if err != nil {
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//store public key on did
	ws.Ctrl.AttachDid(ctx, "did:sonr:temp"+mgoUser.DisplayName, "did:sonr:temp"+credentialID)

	//TODO store cred under user in mgo

	jsonResponse(w, http.StatusText(http.StatusCreated), http.StatusCreated)
}

// GetCredentials gets a user's credentials from the db
func (ws *Server) GetCredentials(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["name"]
	u, err := ws.Ctrl.GetUserByUsername(username)
	if err != nil {
		log.Errorf("user not found: %s: %s", username, err)
		jsonResponse(w, "User not found", http.StatusNotFound)
		return
	}
	cs, err := ws.Ctrl.GetCredentialsForUser(u)
	if err != nil {
		log.Error(err)
		jsonResponse(w, "Credentials not found", http.StatusNotFound)
		return
	}
	jsonResponse(w, cs, http.StatusOK)
}

// DeleteCredential deletes a credential from the db
func (ws *Server) DeleteCredential(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	credID := vars["id"]
	err := ws.Ctrl.DeleteCredentialByID(credID)
	log.Infof("deleting credential: %s", credID)
	if err != nil {
		log.Errorf("error deleting credential: %s", err)
		jsonResponse(w, "Credential not Found", http.StatusNotFound)
		return
	}
	jsonResponse(w, "Success", http.StatusOK)
}

func (ws *Server) HealthHandler(w http.ResponseWriter, req *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

type Response struct {
	Available bool
}

func (ws *Server) CheckName(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	vars := mux.Vars(req)
	name := vars["name"]
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
	did := vars["did"]

	// start := time.Now()
	// e := log.Info()
	// defer func(e *zerolog.Event, start time.Time) {
	// 	if err != nil {
	// 		e = log.Error().Stack().Err(err)
	// 	}
	// 	e.Str("handler", "RegisterName").AnErr("context", ctx.Err()).Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
	// }(e, start)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	//log.Debug().Str("handler", "RegisterName").Bytes("request_body", body).Send()

	var recObj *rt.MsgRegisterName
	err = json.Unmarshal(body, &recObj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	//TODO checkname

	// TODO record name in mongo

	resp, err := ws.Ctrl.RegisterName(ctx, recObj, did)
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
