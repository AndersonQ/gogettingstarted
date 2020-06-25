package solution

import (
	"net/http"

	"github.com/AndersonQ/gogettingstarted/02-http-middlewares/tracking"
)

func TrackingID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := tracking.ContextWithID(r.Context())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
