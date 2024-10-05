package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/middleware"
	"server/views/templates"
)

type AdminHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewAdminHandler(db *sql.DB, slog *slog.Logger) *AdminHandler {
	return &AdminHandler{db: db, slog: slog}
}

func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserCtx(r)
	if userCtx == nil {
		h.slog.Info("User context not found")
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	page := templates.AdminHome()
	err := templates.Admin(page, *userCtx, "Notify-admin").Render(r.Context(), w)
	if err != nil {
		h.slog.Error("failed rendering admin page: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
