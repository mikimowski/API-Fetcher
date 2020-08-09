package handlers

import (
	"context"
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/data"
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDelete(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	req, err := http.NewRequest(http.MethodDelete, "/api/fetcher/1", nil)
	if err != nil {
		t.Fatal(err)
	}
	// populating request content with id
	ctx := req.Context()
	ctx = context.WithValue(ctx, subscriptionIDKey, data.ID(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Run
	handler.Delete(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusNoContent)
	expected := ""
	mock.TestBody(t, rr.Body.String(), expected)
}
