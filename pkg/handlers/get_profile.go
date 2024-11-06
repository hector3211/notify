package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/middleware"
	"server/models"
	"server/pkg/service"
	"server/views/templates"
)

type ProfileHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewProfileHandler(db *sql.DB, slog *slog.Logger) *ProfileHandler {
	return &ProfileHandler{db: db, slog: slog}
}

func (h *ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var isAuth templates.IsAuthenticated
	userCtx := middleware.GetUserCtx(r)
	if userCtx == nil {
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	if userCtx.Role == models.ADMIN {
		w.Header().Set("HX-Redirect", "/admin")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	invoiceService := service.NewInvoiceService(h.db, h.slog)
	jobs := invoiceService.GetUserInvoices(userCtx.ID)

	isAuth = true
	page := templates.Profile(*userCtx, jobs)
	err := templates.Layout(page, isAuth, "profile", "Notify-profile").Render(r.Context(), w)
	if err != nil {
		h.slog.Error("Failed to render profile page: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
