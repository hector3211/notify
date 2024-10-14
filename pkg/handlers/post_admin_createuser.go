package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/models"
	"server/pkg/service"
	"server/views/templates"
	"strings"
)

type PostAdminCreateUserHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewPostAdminCreateUserHandler(db *sql.DB, slog *slog.Logger) *PostAdminCreateUserHandler {
	return &PostAdminCreateUserHandler{db: db, slog: slog}
}

func (h *PostAdminCreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userFirstName := r.FormValue("firstname")
	userLastName := r.FormValue("lastname")
	userEmail := r.FormValue("email")
	userPassword := r.FormValue("password")
	userRole := r.FormValue("role")

	userService := service.NewUserService(h.db)
	if userService.CheckEmailExists(strings.TrimSpace(userEmail)) {
		h.slog.Info("user email already exists")
		err := templates.Toast(models.ErrorNotification, "Email already in use, try another one").Render(r.Context(), w)
		if err != nil {
			h.slog.Error("Failed to Toaster: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusConflict)
		return
	}

	user := &models.User{
		FirstName: userFirstName,
		LastName:  userLastName,
		Password:  strings.TrimSpace(userPassword),
		Email:     strings.TrimSpace(userEmail),
		Role:      models.UserRole(userRole),
	}

	userRes := userService.CreateUser(user)
	if userRes == nil {
		h.slog.Info("user creation failed in admin dashboard")
		err := templates.Toast(models.ErrorNotification, "Oops something went wrong, try again later").Render(r.Context(), w)
		if err != nil {
			h.slog.Error("Failed to Toaster: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := templates.Toast(models.SuccessNotification, "Successfully created user").Render(r.Context(), w)
	if err != nil {
		h.slog.Error("Failed to Toaster: " + err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
}
