package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

// defaultTemplates are included every time a template is rendered.
var defaultTemplates = []string{"./templates/base.html", "./templates/info.html"}

// JSONResponse attempts to set the status code, c, and marshal the given
// interface, d, into a response that is written to the given ResponseWriter.
func jsonResponse(w http.ResponseWriter, d interface{}, c int) {
	dj, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	fmt.Fprintf(w, "%s", dj)
}

// renderTemplate renders the template to the ResponseWriter
func renderTemplate(w http.ResponseWriter, f string, data interface{}) {
	t, err := template.ParseFiles(append(defaultTemplates, fmt.Sprintf("./templates/%s", f))...)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering template: %s", err), http.StatusInternalServerError)
		return
	}
	t.ExecuteTemplate(w, "base", data)
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewEncoder.Encode: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("io.Copy: %v", err)
		return
	}
}
