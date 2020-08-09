package handlers

import (
	"encoding/json"
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/data"
	"net/http"
)

func (s *Subscriptions) Update(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(subscriptionIDKey).(data.ID)

	// get current subscription
	sub, err := s.dao.GetSubscriptionByID(id)
	if err != nil {
		s.l.Infof("unable to update %d: %s", id, err)
		http.Error(w, "unable to update", http.StatusInternalServerError)
		return
	}

	// apply patch - this will overwrite fields specified by user
	err = json.NewDecoder(r.Body).Decode(&sub)
	if err != nil {
		s.l.Infof("invalid json %d: %s", id, err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// validate changes
	// id is compared in case client passed id in json and it was different. If specified, it should match.
	if err = sub.Validate(); err != nil || id != sub.ID {
		s.l.Infof("invalid json %d: %s", id, err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// replace
	_, err = s.dao.ReplaceSubscription(sub)
	if err != nil {
		s.l.Infof("unable to update %d: %s", id, err)
		http.Error(w, "unable to update", http.StatusInternalServerError)
		return
	}
	// Inform subscriber about the update
	s.subscriber.Update(sub)

	s.respondWithJSON(sub, w)
	s.l.Debugf("updated %+v", sub)
}
