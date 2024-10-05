package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/pkg/service"
	"server/views/templates"
	"strconv"
	"strings"
)

type JobHander struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewGetJobHandler(db *sql.DB, slog *slog.Logger) *JobHander {
	return &JobHander{db: db, slog: slog}
}

func (h *JobHander) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userEmail := r.FormValue("email")
	userInvoice := r.FormValue("invoice")
	userService := service.NewUserService(h.db)
	invoiceService := service.NewInvoiceService(h.db, h.slog)

	if !userService.CheckEmailExists(strings.TrimSpace(userEmail)) {
		h.slog.Info("No such email exists in db: " + userEmail)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := userService.GetUserByEmail(userEmail)
	if user == nil {
		h.slog.Info("No such user exists in db: " + userEmail)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	InvoiceIdNum, err := strconv.Atoi(userInvoice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	invoiceId := invoiceService.GetInvoiceId(user.ID, InvoiceIdNum)
	if invoiceId == 0 {
		h.slog.Info("No such invoice exists in db: " + userInvoice)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: fix string to int here
	invoice := invoiceService.GetInvoice(string(invoiceId))
	if invoice == nil {
		h.slog.Info("No such job exists in db with invoice: " + userInvoice)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = templates.Invoice(*invoice).Render(r.Context(), w)
	if err != nil {
		h.slog.Error("Failed to job section: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
