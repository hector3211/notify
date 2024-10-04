package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"server/models"
	"server/pkg/service"
	"server/views/templates"

	"github.com/go-chi/chi/v5"
)

type DeleteAdminJobHandler struct {
	db *sql.DB
}

func NewDeleteAdminJobHandler(db *sql.DB) *DeleteAdminJobHandler {
	return &DeleteAdminJobHandler{
		db: db,
	}
}

func (h *DeleteAdminJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jobID := chi.URLParam(r, "id")
	invoiceService := service.NewInvoiceService(h.db)

	err := invoiceService.DeleteInvoice(jobID)
	if err != nil {
		log.Printf("failed deleting job: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = templates.Toast(models.SuccessNotification, "Successfully deleted job!").Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to Toaster: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
