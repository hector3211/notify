package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/pkg/service"
	"server/views/templates"

	"github.com/go-chi/chi/v5"
)

type GetAdminJobEditHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewGetAdminJobEditHandler(db *sql.DB, slog *slog.Logger) *GetAdminJobEditHandler {
	return &GetAdminJobEditHandler{db: db, slog: slog}
}

func (h *GetAdminJobEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	invID := chi.URLParam(r, "id")
	h.slog.Info("got job ID: " + invID)
	invoiceService := service.NewInvoiceService(h.db, h.slog)
	invoiceData := invoiceService.GetInvoice(invID)
	if invoiceData == nil {
		h.slog.Error("fialed getting invoice data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := templates.AdminEditJobForm(*invoiceData).Render(r.Context(), w)
	if err != nil {
		h.slog.Error("Failed to render admin edit job form page: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
