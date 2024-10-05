package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/middleware"
	"server/pkg/service"
	"server/views/templates"
)

type AdminJobHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewAdminJobHandler(db *sql.DB, slog *slog.Logger) *AdminJobHandler {
	return &AdminJobHandler{db: db, slog: slog}
}

func (h *AdminJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserCtx(r)
	if userCtx == nil {
		slog.Info("User context not found")
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	invoiceService := service.NewInvoiceService(h.db, h.slog)
	jobs := invoiceService.GetAllInvoices()

	page := templates.AdminJobs(jobs)
	err := templates.Admin(page, *userCtx, "Notify-admin").Render(r.Context(), w)
	if err != nil {
		slog.Error("Failed to render admin page: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
