package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ibrahimijai/go-http-server/models"
)

// newTestController wires a controller against the in-memory store.
func newTestController() *URLController {
	logger := slog.New(slog.NewTextHandler(testWriter{}, nil))
	return New(models.NewMemory(), logger)
}

// testWriter swallows log output so test runs stay clean.
type testWriter struct{}

func (testWriter) Write(p []byte) (int, error) { return len(p), nil }

func TestShorten(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "valid url",
			body:       `{"url": "https://go.dev/blog"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "relative url rejected",
			body:       `{"url": "/just/a/path"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "non-http scheme rejected",
			body:       `{"url": "ftp://example.com"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "garbage json",
			body:       `{not json`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "unknown field rejected",
			body:       `{"url": "https://go.dev", "extra": true}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestController()

			req := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			c.Shorten(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}
		})
	}
}

func TestShortenThenRedirect(t *testing.T) {
	c := newTestController()

	req := httptest.NewRequest(http.MethodPost, "/api/shorten",
		strings.NewReader(`{"url": "https://go.dev/blog"}`))
	rec := httptest.NewRecorder()
	c.Shorten(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("shorten status = %d, want 201", rec.Code)
	}

	var resp struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decoding response: %v", err)
	}

	// SetPathValue simulates the router's path parameters when calling the
	// controller directly.
	req = httptest.NewRequest(http.MethodGet, "/"+resp.Code, nil)
	req.SetPathValue("code", resp.Code)
	rec = httptest.NewRecorder()
	c.Redirect(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("redirect status = %d, want 302", rec.Code)
	}
	if loc := rec.Header().Get("Location"); loc != "https://go.dev/blog" {
		t.Errorf("Location = %q, want %q", loc, "https://go.dev/blog")
	}
}

func TestRedirectUnknownCode(t *testing.T) {
	c := newTestController()

	req := httptest.NewRequest(http.MethodGet, "/zzzzzzz", nil)
	req.SetPathValue("code", "zzzzzzz")
	rec := httptest.NewRecorder()
	c.Redirect(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", rec.Code)
	}
}
