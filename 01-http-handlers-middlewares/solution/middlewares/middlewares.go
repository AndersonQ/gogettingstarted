package middlewares

import (
	"net/http"

	"github.com/AndersonQ/gogettingstarted/01-http-handlers-middlewares/tracking"
)

func TrackingID(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := tracking.ContextWithID(r.Context())

		handler(w, r.WithContext(ctx))
	}
}
