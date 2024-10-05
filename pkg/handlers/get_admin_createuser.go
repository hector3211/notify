package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/middleware"
	"server/views/templates"
)

type GetAdminCreateUserHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewGetAdminCreateUserHandler(db *sql.DB, slog *slog.Logger) *GetAdminCreateUserHandler {
	return &GetAdminCreateUserHandler{db: db, slog: slog}
}

func (h *GetAdminCreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserCtx(r)
	if userCtx == nil {
		h.slog.Info("User context not found")
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	page := templates.CreateUserForm()
	err := templates.Admin(page, *userCtx, "Notify-admin").Render(r.Context(), w)
	if err != nil {
		h.slog.Error("Failed to render admin page: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
