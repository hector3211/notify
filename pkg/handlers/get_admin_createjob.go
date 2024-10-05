package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/middleware"
	"server/views/templates"
)

type GetAdminCreateJobHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewGetAdminCreateJobHandler(db *sql.DB, slog *slog.Logger) *GetAdminCreateJobHandler {
	return &GetAdminCreateJobHandler{db: db, slog: slog}
}

func (h *GetAdminCreateJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserCtx(r)
	if userCtx == nil {
		slog.Info("User context not found")
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	page := templates.CreateJobForm()
	err := templates.Admin(page, *userCtx, "Notify-admin").Render(r.Context(), w)
	if err != nil {
		slog.Error("Failed to render admin page: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
