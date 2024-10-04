package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"server/models"
	"server/pkg/service"
	"server/views/templates"
	"strconv"
)

type PostAdminCreateJobHandler struct {
	db *sql.DB
}

func NewPostAdminCreateJobHandler(db *sql.DB) *PostAdminCreateJobHandler {
	return &PostAdminCreateJobHandler{db: db}
}

func (h *PostAdminCreateJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jobFormUserId := r.FormValue("id")
	jobFormInvoiceNumber := r.FormValue("invoice")

	userId, err := strconv.Atoi(jobFormUserId)
	if err != nil {
		http.Error(w, "failed converting userId to int", http.StatusInternalServerError)
		return
	}

	invoiceService := service.NewInvoiceService(h.db)
	err = invoiceService.CreateInvoice(userId, jobFormInvoiceNumber)
	if err != nil {
		http.Error(w, "failed new invoice", http.StatusInternalServerError)
		return
	}

	err = templates.Toast(models.SuccessNotification, "Successfully created job!").Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to Toaster: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
