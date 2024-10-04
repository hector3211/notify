package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"server/pkg/service"
	"server/views/templates"

	"github.com/go-chi/chi/v5"
)

type GetAdminJobEditHandler struct {
	db *sql.DB
}

func NewGetAdminJobEditHandler(db *sql.DB) *GetAdminJobEditHandler {
	return &GetAdminJobEditHandler{db: db}
}

func (h *GetAdminJobEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	invID := chi.URLParam(r, "id")
	log.Printf("got job ID:%s", invID)
	invoiceService := service.NewInvoiceService(h.db)
	invoiceData := invoiceService.GetInvoice(invID)
	if invoiceData == nil {
		log.Println("fialed getting invoice data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := templates.AdminEditJobForm(*invoiceData).Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render admin edit job form page: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
