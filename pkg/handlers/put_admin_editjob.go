package handlers

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"server/middleware"
	"server/models"
	"server/pkg/service"

	"github.com/go-chi/chi/v5"
)

var notifyHostEmail = os.Getenv("NOTIFY_HOST_EMAIL")

type PutAdminJobEditHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewPutAdminJobEditHandler(db *sql.DB, slog *slog.Logger) *PutAdminJobEditHandler {
	return &PutAdminJobEditHandler{db: db, slog: slog}
}

func (h *PutAdminJobEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserCtxFromCookie(w, r)
	if userCtx == nil || userCtx.Role != models.ADMIN {
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	invID := chi.URLParam(r, "id")
	newStatus := r.FormValue("status")
	invoiceService := service.NewInvoiceService(h.db, h.slog)

	err := invoiceService.UpdateInvoiceStatus(models.JobStatus(newStatus), invID)
	if err != nil {
		h.slog.Error(fmt.Sprintf("failed updating inv %s: %v", invID, err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// TODO: get user email
	var userEmail string = ""

	// Send email notification
	if models.JobStatus(newStatus) == models.JOBDONE {
		msg := fmt.Sprintf("Job %s is done!", invID)
		emailService := service.NewEmailService(h.db, h.slog)
		go func() {
			err := emailService.SendEmail(2, userEmail, "Job Status", msg)
			if err != nil {
				h.slog.Error("failed sending email", "err", err.Error())
			}
		}()
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}
