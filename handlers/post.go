package handlers

import (
	"encoding/json"
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/data"
	"net/http"
)

type addSubscriptionResponse struct {
	ID data.ID `json:"id"`
}

func (s *Subscriptions) Add(w http.ResponseWriter, r *http.Request) {
	sub := &data.Subscription{}
	err := json.NewDecoder(r.Body).Decode(sub)
	if err != nil {
		s.l.Infof("invalid json: %s", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if err = sub.Validate(); err != nil {
		s.l.Infof("invalid json: %s", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	sub, err = s.dao.AddSubscription(sub)
	if err != nil {
		s.l.Infof("unable to add subscription: %s", err)
		http.Error(w, "unable to add subscription", http.StatusInternalServerError)
		return
	}
	s.subscriber.Subscribe(sub)

	s.respondWithJSON(addSubscriptionResponse{ID: sub.ID}, w)
	s.l.Debugf("added subscription: %+v", sub)
}
