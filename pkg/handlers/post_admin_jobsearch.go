package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/models"
	"server/pkg/service"
	"server/utils"
	"server/views/templates"
)

type SearchJobHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewPostSearchJobHandler(db *sql.DB, slog *slog.Logger) *SearchJobHandler {
	return &SearchJobHandler{db: db, slog: slog}
}

func (h *SearchJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	invoiceService := service.NewInvoiceService(h.db, h.slog)
	query := r.FormValue("job-query")

	allInvoices := invoiceService.GetAllInvoices()
	var filteredResults []models.Invoice

	for _, res := range allInvoices {
		if query == "" || utils.ContainsQuery(res.Invoice, query) {
			filteredResults = append(filteredResults, res)
		}
	}

	err := templates.JobSearch(filteredResults).Render(r.Context(), w)
	if err != nil {
		h.slog.Error("Failed to render admin page: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
