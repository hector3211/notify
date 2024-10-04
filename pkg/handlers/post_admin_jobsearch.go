package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"server/models"
	"server/pkg/service"
	"server/utils"
	"server/views/templates"
)

type SearchJobHandler struct {
	db *sql.DB
}

func NewPostSearchJobHandler(db *sql.DB) *SearchJobHandler {
	return &SearchJobHandler{db: db}
}

func (h *SearchJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	invoiceService := service.NewInvoiceService(h.db)
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
		log.Printf("Failed to render admin page: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
