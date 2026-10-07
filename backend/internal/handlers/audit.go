package handlers

import (
	"context"
	"log"

	"backend/internal/models"
)

// logAudit writes an audit entry; failures are logged but never abort the request.
func (h *Handler) logAudit(ctx context.Context, entry *models.Log) {
	if _, err := h.app.Models.Log.Insert(ctx, entry); err != nil {
		log.Printf("[AUDIT] failed to write log entry: %v", err)
	}
}
