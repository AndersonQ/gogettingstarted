package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AndersonQ/gogettingstarted/01-http-handlers-middlewares/tracking"
)

func TestTrackingID(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
	w := httptest.NewRecorder()

	var trackingID string
	handler := func(w http.ResponseWriter, r *http.Request) {
		trackingID = tracking.IdFromContext(r.Context())
	}

	h := TrackingID(handler)
	h(w, r)

	if trackingID == "" {
		t.Error("expected a tracking id, got an empty string")
	}
}
