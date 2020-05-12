package handlers

import (
	"fmt"
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Print(fmt.Sprintf("method: %s, URL: %s, headers: %v",
		r.Method,
		r.URL.String(),
		r.Header))

	w.WriteHeader(http.StatusTeapot)
	n, err := w.Write([]byte(`{"hello":"world"}`))

	log.Print(fmt.Sprintf("witten %d bytes, error: %v", n, err))
}
