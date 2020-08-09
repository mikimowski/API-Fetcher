// Naive end to end tests
// Based on chi-router functionalities

package handlers

import (
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestE2EAdd(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()
	chiRouter := getApiFetcherChiRouter(handler)

	reader := strings.NewReader(`{"url": "https://httpbin.org/range/20", "interval":5}`)
	req, err := http.NewRequest(http.MethodPost, "/api/fetcher/", reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Run
	chiRouter.ServeHTTP(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusOK)
	expected := `{"id":3}`
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestE2EAddThenListAll(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()
	chiRouter := getApiFetcherChiRouter(handler)

	reader := strings.NewReader(`{"url": "https://httpbin.org/range/20", "interval":5}`)
	req, err := http.NewRequest(http.MethodPost, "/api/fetcher/", reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Run
	chiRouter.ServeHTTP(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusOK)
	expected := `{"id":3}`
	mock.TestBody(t, rr.Body.String(), expected)

	// ListAll
	// Setup
	req, err = http.NewRequest(http.MethodGet, "/api/fetcher", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	// Run
	chiRouter.ServeHTTP(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusOK)
	expected = `[{"id":1,"url":"https://httpbin.org/range/15","interval":60},{"id":2,"url":"https://httpbin.org/delay/10","interval":120},{"id":3,"url":"https://httpbin.org/range/20","interval":5}]`
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestE2EDeleteThenGetHistory(t *testing.T) {
	// DELETE
	// Setup
	handler := getHandlerMockMemoryDB()
	chiRouter := getApiFetcherChiRouter(handler)

	req, err := http.NewRequest(http.MethodDelete, "/api/fetcher/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Run
	chiRouter.ServeHTTP(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusNoContent)
	expected := ""
	mock.TestBody(t, rr.Body.String(), expected)

	// GET
	// Setup
	req, err = http.NewRequest(http.MethodGet, "/api/fetcher/1/history", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	// Run
	chiRouter.ServeHTTP(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusNotFound)
	expected = "404 page not found\n"
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestE2EDeleteThenListAll(t *testing.T) {
	// DELETE
	// Setup
	handler := getHandlerMockMemoryDB()
	chiRouter := getApiFetcherChiRouter(handler)

	req, err := http.NewRequest(http.MethodDelete, "/api/fetcher/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Run
	chiRouter.ServeHTTP(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusNoContent)
	expected := ""
	mock.TestBody(t, rr.Body.String(), expected)

	// GET
	// Setup
	req, err = http.NewRequest(http.MethodGet, "/api/fetcher", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	// Run
	chiRouter.ServeHTTP(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusOK)
	expected = `[{"id":2,"url":"https://httpbin.org/delay/10","interval":120}]`
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestE2EListHistory(t *testing.T) {
	// GET
	// Setup
	handler := getHandlerMockMemoryDB()
	chiRouter := getApiFetcherChiRouter(handler)

	req, err := http.NewRequest(http.MethodGet, "/api/fetcher/1/history", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Run
	chiRouter.ServeHTTP(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusOK)
	expected := `[{"response":"my mock history","duration":0.532,"created_at":"1559034938.638"},{"response":null,"duration":5,"created_at":"1559034938.638"}]`
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestE2EListHistoryInvalidID(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()
	chiRouter := getApiFetcherChiRouter(handler)

	req, err := http.NewRequest(http.MethodGet, "/api/fetcher/abcd/history", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Run
	chiRouter.ServeHTTP(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusNotFound)
	expected := "404 page not found\n"
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestE2EListHistoryNonExistingID(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()
	chiRouter := getApiFetcherChiRouter(handler)

	req, err := http.NewRequest(http.MethodGet, "/api/fetcher/99/history", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Run
	chiRouter.ServeHTTP(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusNotFound)
	expected := "404 page not found\n"
	mock.TestBody(t, rr.Body.String(), expected)
}
