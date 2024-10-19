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

type AdminAccountHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewAdminAccountHanlder(db *sql.DB, slog *slog.Logger) *AdminAccountHandler {
	return &AdminAccountHandler{
		db:   db,
		slog: slog,
	}
}

func (h *AdminAccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.GetUserCtx(r)
	if userCtx == nil || userCtx.Role != models.ADMIN {
		h.slog.Info("User context not found")
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userService := service.NewUserService(h.db)
	userResponse := userService.GetUserByID(userCtx.ID)

	page := templates.AdminAccount(*userResponse)
	err := templates.Admin(page, *userCtx, "Notify-admin").Render(r.Context(), w)
	if err != nil {
		h.slog.Error("failed rendering admin account page: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
