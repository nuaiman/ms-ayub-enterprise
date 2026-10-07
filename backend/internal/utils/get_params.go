package utils

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GetParamID extracts the "id" URL parameter using chi's router.
// Chi does NOT populate r.PathValue() — that's only for net/http.ServeMux (Go 1.22+).
func GetParamID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ErrorJson(w, http.StatusBadRequest, "invalid parameter")
		return 0, false
	}

	return id, true
}

// GetURLParam returns a named URL parameter from chi.
func GetURLParam(r *http.Request, key string) string {
	return chi.URLParam(r, key)
}

// ParseInt64 converts a string to int64
func ParseInt64(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}
