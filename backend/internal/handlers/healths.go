package handlers

import (
	"backend/internal/utils"
	"net/http"
)

func (h *Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	utils.SuccessJson(w, http.StatusOK, "API is active", nil)
}
