package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/models"
	"server/pkg/service"
	"server/views/templates"
	"strconv"
)

type PostAdminCreateJobHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewPostAdminCreateJobHandler(db *sql.DB, slog *slog.Logger) *PostAdminCreateJobHandler {
	return &PostAdminCreateJobHandler{db: db, slog: slog}
}

func (h *PostAdminCreateJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	jobFormUserId := r.FormValue("id")
	jobFormInvoiceNumber := r.FormValue("invoice")

	userId, err := strconv.Atoi(jobFormUserId)
	if err != nil {
		http.Error(w, "failed converting userId to int", http.StatusInternalServerError)
		return
	}

	invoiceService := service.NewInvoiceService(h.db, h.slog)
	err = invoiceService.CreateInvoice(userId, jobFormInvoiceNumber)
	if err != nil {
		http.Error(w, "failed new invoice", http.StatusInternalServerError)
		return
	}

	err = templates.Toast(models.SuccessNotification, "Successfully created job!").Render(r.Context(), w)
	if err != nil {
		h.slog.Error("Failed to Toaster: " + err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}
