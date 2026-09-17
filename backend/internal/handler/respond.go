package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"travelcrm/internal/repository"
)

// writeJSON encodes v as JSON with the given status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError maps a domain error to an HTTP status + readable message (§14).
func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		writeJSON(w, http.StatusNotFound, errorBody{Error: "Not found"})
	case errors.Is(err, repository.ErrSeatTaken):
		writeJSON(w, http.StatusConflict, errorBody{Error: "Selected seat is no longer available"})
	case errors.Is(err, repository.ErrInvalidInput):
		writeJSON(w, http.StatusBadRequest, errorBody{Error: err.Error()})
	default:
		writeJSON(w, http.StatusInternalServerError, errorBody{Error: "internal server error"})
	}
}

type errorBody struct {
	Error string `json:"error"`
}

// decodeJSON reads and decodes a JSON request body into dst. The underlying
// parse error is wrapped so the client gets a useful 400 message instead of a
// bare "invalid input".
func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("%w: %v", repository.ErrInvalidInput, err)
	}
	return nil
}
