package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/models"
	"server/pkg/jwt"
	"server/pkg/service"
	"server/utils"
	"server/views/templates"
	"strings"
	"time"
)

type PostLoginHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewPostLoginHandler(db *sql.DB, slog *slog.Logger) *PostLoginHandler {
	return &PostLoginHandler{db: db, slog: slog}
}

func (h *PostLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logInEmail := r.FormValue("email")
	logInPassword := r.FormValue("password")
	userService := service.NewUserService(h.db)

	// h.slog.Info("login Password: %s\n", logInPassword)

	if !userService.CheckEmailExists(strings.TrimSpace(logInEmail)) {
		h.slog.Info("No such email exists in db: " + logInEmail)
		err := templates.Toast(models.ErrorNotification, "No user found with matching email").Render(r.Context(), w)
		if err != nil {
			h.slog.Error("Failed to Toaster: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userRes := userService.GetUserByEmail(strings.TrimSpace(logInEmail))
	if userRes == nil {
		h.slog.Info("No User found with email: " + logInEmail)
		err := templates.Toast(models.ErrorNotification, "Oops something went wrong, try again later").Render(r.Context(), w)
		if err != nil {
			h.slog.Error("Failed to Toaster: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !userService.CheckPasswordHash(strings.TrimSpace(logInPassword), userRes.ID) {
		h.slog.Info("passwords did NOT match with hash")
		err := templates.Toast(models.ErrorNotification, "Invalid password no user found").Render(r.Context(), w)
		if err != nil {
			h.slog.Error("Failed to Toaster: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := jwt.NewJwtService().Init(userRes.ID, userRes.Role.String())
	if err != nil {
		h.slog.Info("failed creating jwt when logging in")
		err := templates.Toast(models.ErrorNotification, "Oops something went wrong, try again later").Render(r.Context(), w)
		if err != nil {
			h.slog.Error("Failed to Toaster: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 72),
		HttpOnly: true,
		Secure:   utils.IsProduction(),
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	w.Header().Set("HX-Redirect", "/profile")
	w.WriteHeader(http.StatusSeeOther)
}
