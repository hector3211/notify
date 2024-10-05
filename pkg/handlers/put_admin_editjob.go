package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/models"
	"server/pkg/service"

	"github.com/go-chi/chi/v5"
)

type PutAdminJobEditHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewPutAdminJobEditHandler(db *sql.DB, slog *slog.Logger) *PutAdminJobEditHandler {
	return &PutAdminJobEditHandler{db: db, slog: slog}
}

func (h *PutAdminJobEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	invID := chi.URLParam(r, "id")
	newStatus := r.FormValue("status")
	invoiceService := service.NewInvoiceService(h.db, h.slog)

	err := invoiceService.UpdateInvoiceStatus(models.JobStatus(newStatus), invID)
	if err != nil {
		h.slog.Error("failed updating inv ID: %s with err: %v", invID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}
