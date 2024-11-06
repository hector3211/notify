package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/middleware"
	"server/models"
	"server/views/templates"

	"github.com/go-chi/chi/v5"
)

type GetAdminCreateJobHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewGetAdminCreateJobHandler(db *sql.DB, slog *slog.Logger) *GetAdminCreateJobHandler {
	return &GetAdminCreateJobHandler{db: db, slog: slog}
}

func (h *GetAdminCreateJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paramUserId := chi.URLParam(r, "userid")
	userCtx := middleware.GetUserCtxFromCookie(w, r)
	if userCtx == nil {
		h.slog.Info("User context not found")
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if userCtx.Role != models.ADMIN {
		h.slog.Info("user unauthorized")
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	if paramUserId == "" {
		page := templates.CreateJobForm("")
		err := templates.Admin(page, *userCtx, "Notify-admin").Render(r.Context(), w)
		if err != nil {
			slog.Error("Failed to render admin page: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	page := templates.CreateJobForm(paramUserId)
	err := templates.Admin(page, *userCtx, "Notify-admin").Render(r.Context(), w)
	if err != nil {
		slog.Error("Failed to render admin page: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
