package handlers

import (
	"encoding/json"
	"net/http"
)

// Encodes given interface to json and writes it to http.ResponseWriter.
// No further writing!
func (s *Subscriptions) respondWithJSON(v interface{}, w http.ResponseWriter) {
	e, err := json.Marshal(v)
	if err != nil {
		s.l.Infof("unable to encode data to json: %s", err)
		http.Error(w, "unable to encode data to json", http.StatusInternalServerError)
		return
	}
	if _, err = w.Write(e); err != nil {
		s.l.Infof("failed responding to client", err)
	}
}
