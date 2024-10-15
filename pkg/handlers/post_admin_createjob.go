package handlers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"server/models"
	"server/pkg/service"
	"server/views/templates"
	"strconv"
	"time"
)

type PostAdminCreateJobHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewPostAdminCreateJobHandler(db *sql.DB, slog *slog.Logger) *PostAdminCreateJobHandler {
	return &PostAdminCreateJobHandler{db: db, slog: slog}
}

func (h *PostAdminCreateJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	formUserId := r.FormValue("id")
	formInvoiceNumber := r.FormValue("invoice")
	formInstallDate := r.FormValue("install_date")

	if formInstallDate == "" {
		h.slog.Error("installion date empty")
		err := templates.Toast(models.ErrorNotification, "No installation date provided").Render(r.Context(), w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.slog.Info(fmt.Sprintf("Date: %s", formInstallDate))

	parsedDate, err := time.Parse("2006-01-02", formInstallDate)
	if err != nil {
		h.slog.Error(fmt.Sprintf("failed parsing date time %v", err))
		err = templates.Toast(models.ErrorNotification, "Oops something went wrong, try again later").Render(r.Context(), w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userId, err := strconv.Atoi(formUserId)
	if err != nil {
		h.slog.Error("failed strconv on userid")
		err = templates.Toast(models.ErrorNotification, "Oops something went wrong, try again later").Render(r.Context(), w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	invoiceService := service.NewInvoiceService(h.db, h.slog)
	err = invoiceService.CreateInvoice(userId, formInvoiceNumber, parsedDate.Format("01-02-2006"))
	if err != nil {
		h.slog.Error("failed creating new invoice: " + err.Error())
		err = templates.Toast(models.ErrorNotification, "Oops something went wrong, try again later").Render(r.Context(), w)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = templates.Toast(models.SuccessNotification, "Successfully created job!").Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
