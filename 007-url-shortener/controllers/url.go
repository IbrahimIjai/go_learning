package controllers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/ibrahimijai/go-http-server/helpers"
	"github.com/ibrahimijai/go-http-server/models"
)

// URLController depends on the URLStore interface, not a concrete store,
// which keeps it swappable and testable.
type URLController struct {
	store  models.URLStore
	logger *slog.Logger
}

func New(st models.URLStore, logger *slog.Logger) *URLController {
	return &URLController{store: st, logger: logger}
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// Shorten handles POST /api/shorten
func (c *URLController) Shorten(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req shortenRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body"})
		return
	}

	u, err := url.Parse(req.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		helpers.WriteJSON(w, http.StatusBadRequest, errorResponse{Error: "url must be an absolute http(s) URL"})
		return
	}

	// Generate a code; retry a few times on the rare collision.
	var code string
	for attempt := 0; attempt < 3; attempt++ {
		code, err = helpers.NewCode()
		if err != nil {
			helpers.WriteJSON(w, http.StatusInternalServerError, errorResponse{Error: "could not generate code"})
			return
		}

		err = c.store.Save(r.Context(), code, req.URL)
		if err == nil {
			break
		}
		if !errors.Is(err, models.ErrConflict) {
			c.logger.Error("save failed", "err", err)
			helpers.WriteJSON(w, http.StatusInternalServerError, errorResponse{Error: "storage error"})
			return
		}
	}
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, errorResponse{Error: "could not allocate code"})
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, shortenResponse{
		Code:     code,
		ShortURL: "http://" + r.Host + "/" + code,
	})
}

// Redirect handles GET /{code}
func (c *URLController) Redirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")

	if !helpers.Valid(code) {
		http.NotFound(w, r)
		return
	}

	target, err := c.store.GetByCode(r.Context(), code)
	switch {
	case errors.Is(err, models.ErrNotFound):
		http.NotFound(w, r)
		return
	case err != nil:
		c.logger.Error("lookup failed", "err", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// 302 (Found), not 301: permanent redirects get cached aggressively by
	// browsers, which would make click analytics impossible later.
	http.Redirect(w, r, target, http.StatusFound)
}

// Health handles GET /healthz
func (c *URLController) Health(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
