package handlers

import (
	"github.com/mikimowski/TWFjaWVqLU1pa3XFgmE/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	reader := strings.NewReader(`{"url": "https://httpbin.org/range/", "interval":5}`)
	req, err := http.NewRequest("POST", "/api/fetcher/", reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Run
	handler.Add(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusOK)
	expected := `{"id":3}`
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestAddMissingURL(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	reader := strings.NewReader(`{"interval":5}`)
	req, err := http.NewRequest("POST", "/api/fetcher/", reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Run
	handler.Add(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusBadRequest)
	expected := "invalid json\n" // http.Error uses fmt.FPrintln
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestAddMissingInterval(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	reader := strings.NewReader(`{"url": "https://httpbin.org/range/"}`)
	req, err := http.NewRequest("POST", "/api/fetcher/", reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Run
	handler.Add(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusBadRequest)
	expected := "invalid json\n" // http.Error uses fmt.FPrintln
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestAddNullURL(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	reader := strings.NewReader(`{"url": null, "interval":5}`)
	req, err := http.NewRequest("POST", "/api/fetcher/", reader)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Run
	handler.Add(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusBadRequest)
	expected := "invalid json\n" // http.Error uses fmt.FPrintln
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestAddNullInterval(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	reader := strings.NewReader(`{"url": "https://httpbin.org/range/", "interval":null}`)
	req, err := http.NewRequest("POST", "/api/fetcher/", reader)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	// Run
	handler.Add(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusBadRequest)
	expected := "invalid json\n" // http.Error uses fmt.FPrintln
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestAddNegativeInterval(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	reader := strings.NewReader(`{"url": "https://httpbin.org/range/", "interval":-5}`)
	req, err := http.NewRequest("POST", "/api/fetcher/", reader)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Run
	handler.Add(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusBadRequest)
	expected := "invalid json\n" // http.Error uses fmt.FPrintln
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestAddZeroInterval(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	reader := strings.NewReader(`{"url": "https://httpbin.org/range/", "interval":0}`)
	req, err := http.NewRequest("POST", "/api/fetcher/", reader)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Run
	handler.Add(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusBadRequest)
	expected := "invalid json\n" // http.Error uses fmt.FPrintln
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestAddDecimalInterval(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	reader := strings.NewReader(`{"url": "https://httpbin.org/range/", "interval":3.14}`)
	req, err := http.NewRequest("POST", "/api/fetcher/", reader)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Run
	handler.Add(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusBadRequest)
	expected := "invalid json\n" // http.Error uses fmt.FPrintln
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestAddStringInterval(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	reader := strings.NewReader(`{"url": "https://httpbin.org/range/", "interval":"5"}`)
	req, err := http.NewRequest("POST", "/api/fetcher/", reader)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Run
	handler.Add(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusBadRequest)
	expected := "invalid json\n" // http.Error uses fmt.FPrintln
	mock.TestBody(t, rr.Body.String(), expected)
}

func TestAddNumericalURL(t *testing.T) {
	// Setup
	handler := getHandlerMockMemoryDB()

	reader := strings.NewReader(`{"url": 42, "interval":5}`)
	req, err := http.NewRequest("POST", "/api/fetcher/", reader)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Run
	handler.Add(rr, req)

	// Check
	mock.TestStatus(t, rr.Code, http.StatusBadRequest)
	expected := "invalid json\n" // http.Error uses fmt.FPrintln
	mock.TestBody(t, rr.Body.String(), expected)
}
