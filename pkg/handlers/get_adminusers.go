package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/middleware"
	"server/pkg/service"
	"server/views/templates"
)

type AdminUserHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewAdminUserHandler(db *sql.DB, logger *slog.Logger) *AdminUserHandler {
	return &AdminUserHandler{db: db, slog: logger}
}

func (h *AdminUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserCtx(r)
	if userCtx == nil {
		h.slog.Info("user ctx not found")
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userService := service.NewUserService(h.db)
	users := userService.GetAllUsers()

	page := templates.AdminUsers(*userCtx, users)
	err := templates.Admin(page, *userCtx, "Notify-admin").Render(r.Context(), w)
	if err != nil {
		h.slog.Error("failed rendering admin users page: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
