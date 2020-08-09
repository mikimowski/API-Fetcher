package handlers

import (
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/data"
	"net/http"
)

// Lists all subscriptions
func (s *Subscriptions) ListAll(w http.ResponseWriter, r *http.Request) {
	subs, err := s.dao.GetAllSubscriptions()
	if err != nil {
		s.l.Infof("unable to fetch subscription data: %s", err)
		http.Error(w, "unable to fetch subscription data", http.StatusInternalServerError)
		return
	}

	s.respondWithJSON(subs, w)
	s.l.Debugf("subscriptions listAll: %+v", subs)
}

// Lists history for subscription in context
func (s *Subscriptions) ListHistory(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(subscriptionIDKey).(data.ID)
	history, err := s.dao.GetSubscriptionHistoryByID(id)
	if err != nil {
		s.l.Infof("unable to fetch history data: %s", err)
		http.Error(w, "unable to fetch history data", http.StatusInternalServerError)
		return
	}

	s.respondWithJSON(history, w)
	s.l.Debugf("subscriptions listHistory: %+v", history)
}
