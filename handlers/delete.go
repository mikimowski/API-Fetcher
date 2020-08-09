package handlers

import (
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/data"
	"net/http"
)

// Deletes subscription with id from context
func (s *Subscriptions) Delete(w http.ResponseWriter, r *http.Request) {
	// Parsed and validated via middleware
	id := r.Context().Value(subscriptionIDKey).(data.ID)

	// IMMEDIATELY delete subscription so that following calls to API reflect it correctly
	if _, err := s.dao.DeleteSubscription(id); err != nil {
		s.l.Infof("unable to delete subscription: %s", err)
		http.Error(w, "unable to delete subscription", http.StatusInternalServerError)
		return
	}

	// unsubscribe IMMEDIATELY,
	// actual clean up of the history can be done later -> this goroutine does not bother,
	// it only informs Subscriber about it and respond to client ASAP (see subscriber implementation for details)
	s.subscriber.Unsubscribe(id)

	w.WriteHeader(http.StatusNoContent)
	s.l.Infof("subscriptions Delete: %d", id)
}
