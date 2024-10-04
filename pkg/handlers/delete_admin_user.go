package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"server/models"
	"server/pkg/service"
	"server/views/templates"

	"github.com/go-chi/chi/v5"
)

type DeleteAdminUserHandler struct {
	db *sql.DB
}

func NewDeleteAdminUserHandler(db *sql.DB) *DeleteAdminUserHandler {
	return &DeleteAdminUserHandler{
		db: db,
	}
}

func (h *DeleteAdminUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	userService := service.NewUserService(h.db)

	err := userService.DeleteUser(userID)
	if err != nil {
		log.Printf("failed deleting user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = templates.Toast(models.SuccessNotification, "Successfully deleted user!").Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to Toaster: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
