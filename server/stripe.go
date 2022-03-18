package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sonr-io/webauthn.io/models"
)

func (ws *Server) CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	req := models.SnrItem{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return
	}

	pi, err := ws.Ctrl.StripeIntent(req)
	log.Printf("pi.New: %v", pi.ClientSecret)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("pi.New: %v", err)
		return
	}

	writeJSON(w, struct {
		ClientSecret string `json:"clientSecret"`
	}{
		ClientSecret: pi.ClientSecret,
	})
}
