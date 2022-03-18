package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sonr-io/webauthn.io/models"
)

func (ws *Server) CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Items []models.SnrItem `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewDecoder.Decode: %v", err)
		return
	}

	if len(req.Items) != 1 {
		//throw error saying not enough items
		return
	}

	pi, err := ws.Ctrl.StripeIntent(req.Items[0])
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
