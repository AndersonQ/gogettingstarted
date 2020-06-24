package middlewares

import (
	"net/http"

	"github.com/rs/zerolog"
)

func TrackingID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func LogRequest(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			logger.Info().Msg("request finished")
		})
	}

}

func a() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	lm := LogRequest(zerolog.Logger{})
	hh := lm(h)

	hh.ServeHTTP(nil, nil)
}
