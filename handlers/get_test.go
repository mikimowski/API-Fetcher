package handlers

import (
	"context"
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/data"
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListAll(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	req, err := http.NewRequest(http.MethodGet, "/api/fetcher/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Run
	handler.ListAll(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusOK)
	// Actually... ordering is not guaranteed, therefore this kind of testing might fail just because of ordering :)
	expected := `[{"id":1,"url":"https://httpbin.org/range/15","interval":60},{"id":2,"url":"https://httpbin.org/delay/10","interval":120}]`
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestListAllDatabaseError(t *testing.T) {
	// Setup
	handler := getHandlerFailingDB()

	req, err := http.NewRequest(http.MethodGet, "/api/fetcher/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Run
	handler.ListAll(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusInternalServerError)
}

func TestListHistory(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	req, err := http.NewRequest(http.MethodGet, "/api/fetcher/1/history", nil)
	if err != nil {
		t.Fatal(err)
	}
	// populating request content with id
	ctx := req.Context()
	ctx = context.WithValue(ctx, subscriptionIDKey, data.ID(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Run
	handler.ListHistory(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusOK)
	expected := `[{"response":"my mock history","duration":0.532,"created_at":"1559034938.638"},{"response":null,"duration":5,"created_at":"1559034938.638"}]`
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestListHistoryDatabaseError(t *testing.T) {
	// Setup
	handler := getHandlerFailingDB()

	req, err := http.NewRequest(http.MethodGet, "/api/fetcher/1/history", nil)
	if err != nil {
		t.Fatal(err)
	}
	// populating request content with id
	ctx := req.Context()
	ctx = context.WithValue(ctx, subscriptionIDKey, data.ID(1))
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Run
	handler.ListHistory(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusInternalServerError)
}
