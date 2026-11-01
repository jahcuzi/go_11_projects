package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleFormPost(t *testing.T) {
	body := strings.NewReader("name=Aman&address=Hyderabad")
	req := httptest.NewRequest(http.MethodPost, "/form", body)
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	handleForm(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("Got status %d, want %d", rec.Code, http.StatusOK)
	}

	for _, want := range []string{"POST request successful", "name is Aman", "Address is Hyderabad"} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Errorf("response body %q doesn't contain %q", rec.Body.String(), want)
		}
	}
}
