package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/models"
	"server/pkg/service"
	"server/views/templates"

	"github.com/go-chi/chi/v5"
)

type DeleteAdminJobHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewDeleteAdminJobHandler(db *sql.DB, slog *slog.Logger) *DeleteAdminJobHandler {
	return &DeleteAdminJobHandler{db: db, slog: slog}
}

func (h *DeleteAdminJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")
	invoiceService := service.NewInvoiceService(h.db, h.slog)

	err := invoiceService.DeleteInvoice(jobID)
	if err != nil {
		h.slog.Error("failed deleteing job: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = templates.Toast(models.SuccessNotification, "Successfully deleted job!").Render(r.Context(), w)
	if err != nil {
		h.slog.Error("failed to toaster up: " + err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
