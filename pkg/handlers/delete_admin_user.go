package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/models"
	"server/pkg/service"
	"server/views/templates"

	"github.com/go-chi/chi/v5"
)

type DeleteAdminUserHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewDeleteAdminUserHandler(db *sql.DB, slog *slog.Logger) *DeleteAdminUserHandler {
	return &DeleteAdminUserHandler{db: db, slog: slog}
}

func (h *DeleteAdminUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	userService := service.NewUserService(h.db)

	err := userService.DeleteUser(userID)
	if err != nil {
		h.slog.Error("failed deleting user: " + err.Error())
		err = templates.Toast(models.ErrorNotification, "Oops something went wrong, try again later").Render(r.Context(), w)
		if err != nil {
			h.slog.Error("failed to toaster up: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = templates.Toast(models.SuccessNotification, "Successfully deleted user").Render(r.Context(), w)
	if err != nil {
		h.slog.Error("failed to toaster up: " + err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
