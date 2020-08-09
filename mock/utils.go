// Utility functions for testing

package mock

import "testing"

func TestStatus(t *testing.T, received, expected int) {
	if expected != received {
		t.Errorf("wrong status code: got '%v' expected '%v'", received, expected)
	}
}

func TestBody(t *testing.T, received, expected string) {
	if expected != received {
		t.Errorf("unexpected body: got '%v' expected '%v'", received, expected)
	}
}
