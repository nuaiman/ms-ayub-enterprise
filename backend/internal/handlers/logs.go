package handlers

import (
	"backend/internal/models"
	"backend/internal/utils"
	"log"
	"net/http"
)

func (h *Handler) GetAllLogsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[LOGS] GetAllLogsHandler called - Method: %s, Path: %s", r.Method, r.URL.Path)

	logs, err := h.app.Models.Log.GetAll(r.Context())
	if err != nil {
		log.Printf("[LOGS] GetAllLogsHandler ERROR: %v", err)
		utils.ErrorJson(w, http.StatusInternalServerError, "failed to fetch logs")
		return
	}

	if logs == nil {
		logs = []models.Log{}
	}

	log.Printf("[LOGS] GetAllLogsHandler SUCCESS: fetched %d logs", len(logs))
	utils.SuccessJson(w, http.StatusOK, "logs fetched successfully", logs)
}
