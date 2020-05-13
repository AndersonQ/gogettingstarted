package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	expectedResponseCode := http.StatusTeapot
	expectedBody := `{"hello":"world"}`

	buff := &bytes.Buffer{}
	r := httptest.NewRequest(http.MethodGet, "https://go.lang/", nil)
	r.Header.Set("x-mascot", "Go Gopher")
	w := httptest.NewRecorder()
	w.Body = buff

	Handler(w, r)

	responseCode := w.Code
	body := w.Body.String()

	if responseCode != expectedResponseCode {
		t.Errorf("want: %d-%s, got: %d-%s",
			expectedResponseCode,
			http.StatusText(expectedResponseCode),
			responseCode,
			http.StatusText(responseCode))
	}
	if body != expectedBody {
		t.Errorf("want: %s, got: %s", expectedBody, body)
	}
}
