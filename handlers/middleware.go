package handlers

import (
	"bytes"
	"context"
	"github.com/go-chi/chi"
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/data"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Below does not apply in our short program but it's a good practice ;)
// It prevents collisions with context keys from different packages
// The key type is unexported to prevent collisions with context keys defined in other packages.
type idKey data.ID

// Context key for subscription ID
const subscriptionIDKey idKey = 0

// Middleware to load id into context and validate it.
// id is valid if subscription with given id is stored in database,
// otherwise HTTP 404 (Not Found) is returned to the client.
func (s *Subscriptions) IDContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idURL := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idURL)
		if err != nil {
			s.l.Infof("%s", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		exists, err := s.dao.Exists(data.ID(id))
		if err != nil {
			s.l.Infof("%s", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if exists {
			ctx := context.WithValue(r.Context(), subscriptionIDKey, data.ID(id))
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		} else {
			http.NotFound(w, r)
		}
	})
}

// Limits payload size to payloadLimit
func (s *Subscriptions) PayloadLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, payloadLimit)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			s.l.Info(err.Error())
			switch {
			case err.Error() == "http: request body too large":
				http.Error(w, "payload limit exceeded", http.StatusRequestEntityTooLarge)
				return
			default:
				http.Error(w, "error reading request", http.StatusInternalServerError)
				return
			}
		}

		// TODO (to discuss during review) dilemma: reading body twice vs inserting body slice to context
		//ctx := context.WithValue(r.Context(),"xd", body)
		//r = r.WithContext(ctx)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
		return
	})
}

// Sets content type to "application/json"
func (s *Subscriptions) ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
