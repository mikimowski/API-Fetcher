package handlers

import (
	"context"
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/data"
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdate(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	reader := strings.NewReader(`{"url": "https://httpbin.org/range/10", "interval":13}`)
	req, err := http.NewRequest(http.MethodPatch, "/api/fetcher/1", reader)
	if err != nil {
		t.Fatal(err)
	}

	// populating request content with id
	ctx := req.Context()
	ctx = context.WithValue(ctx, subscriptionIDKey, data.ID(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Run
	handler.Update(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusOK)
	expected := `{"id":1,"url":"https://httpbin.org/range/10","interval":13}`
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestUpdateURL(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	reader := strings.NewReader(`{"url": "https://httpbin.org/range/10"}`)
	req, err := http.NewRequest(http.MethodPatch, "/api/fetcher/1", reader)
	if err != nil {
		t.Fatal(err)
	}

	// populating request content with id
	ctx := req.Context()
	ctx = context.WithValue(ctx, subscriptionIDKey, data.ID(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Run
	handler.Update(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusOK)
	expected := `{"id":1,"url":"https://httpbin.org/range/10","interval":60}`
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestUpdateInterval(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	reader := strings.NewReader(`{"interval":5}`)
	req, err := http.NewRequest(http.MethodPatch, "/api/fetcher/1", reader)
	if err != nil {
		t.Fatal(err)
	}

	// populating request content with id
	ctx := req.Context()
	ctx = context.WithValue(ctx, subscriptionIDKey, data.ID(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Run
	handler.Update(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusOK)
	expected := `{"id":1,"url":"https://httpbin.org/range/15","interval":5}`
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestUpdateIntervalInvalid(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	reader := strings.NewReader(`{"interval":-13}`)
	req, err := http.NewRequest(http.MethodPatch, "/api/fetcher/1", reader)
	if err != nil {
		t.Fatal(err)
	}

	// populating request content with id
	ctx := req.Context()
	ctx = context.WithValue(ctx, subscriptionIDKey, data.ID(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Run
	handler.Update(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusBadRequest)
	expected := "invalid json\n"
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestUpdateWrongID(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	reader := strings.NewReader(`{"id":13,"url":"https://httpbin.org/range/15","interval":5}`)
	req, err := http.NewRequest(http.MethodPatch, "/api/fetcher/1", reader)
	if err != nil {
		t.Fatal(err)
	}

	// populating request content with id
	ctx := req.Context()
	ctx = context.WithValue(ctx, subscriptionIDKey, data.ID(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Run
	handler.Update(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusBadRequest)
	expected := "invalid json\n"
	mock.TestBody(t, rr.Body.String(), expected)
}
