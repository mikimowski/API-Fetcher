package handlers

import (
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func generateSeq(n int64) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = 'a'
	}
	return string(b)
}

// Send body of size == payloadLimit + 1
func TestPayloadLimitExceeded(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("payload exceeded and middleware didn't react")
	})

	reader := strings.NewReader(generateSeq(payloadLimit + 1))
	req, err := http.NewRequest("POST", "/api/fetcher/", reader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Run
	handler.PayloadLimit(nextHandler).ServeHTTP(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusRequestEntityTooLarge)
}

// Send body of size == payloadLimit
func TestPayloadLimitNotExceeded(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	reader := strings.NewReader(generateSeq(payloadLimit))
	req, err := http.NewRequest("POST", "/api/fetcher/", reader)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Run
	handler.PayloadLimit(nextHandler).ServeHTTP(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusOK)
}
