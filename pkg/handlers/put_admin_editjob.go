package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"server/models"
	"server/pkg/service"

	"github.com/go-chi/chi/v5"
)

type PutAdminJobEditHandler struct {
	db *sql.DB
}

func NewPutAdminJobEditHandler(db *sql.DB) *PutAdminJobEditHandler {
	return &PutAdminJobEditHandler{db: db}
}

func (h *PutAdminJobEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	invID := chi.URLParam(r, "id")
	newStatus := r.FormValue("status")
	invoiceService := service.NewInvoiceService(h.db)

	err := invoiceService.UpdateInvoiceStatus(models.JobStatus(newStatus), invID)
	if err != nil {
		log.Printf("failed updating inv ID: %s with err: %v", invID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}
