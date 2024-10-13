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

type PostSignUpHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewPostSignupHandler(db *sql.DB, slog *slog.Logger) *PostSignUpHandler {
	return &PostSignUpHandler{db: db, slog: slog}
}

func (h *PostSignUpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	signupFirstName := r.FormValue("firstname")
	signupLastName := r.FormValue("lastname")
	signupEmail := r.FormValue("email")
	signupPassword := r.FormValue("password")

	// h.slog.Info("signup email: %s\n", signupEmail)
	// h.slog.Info("signup password: %s\n", signupPassword)

	userService := service.NewUserService(h.db)
	if userService.CheckEmailExists(strings.TrimSpace(signupEmail)) {
		h.slog.Info("user with email already exists")
		err := templates.Toast(models.ErrorNotification, "Email already in use, try another one").Render(r.Context(), w)
		if err != nil {
			h.slog.Error("Failed to Toaster: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newUser := &models.User{
		FirstName: signupFirstName,
		LastName:  signupLastName,
		Password:  signupPassword,
		Email:     signupEmail,
		Role:      models.UserRole(models.USER),
	}

	userRes := userService.CreateUser(newUser)
	if userRes == nil {
		h.slog.Info("failed creating new user")
		err := templates.Toast(models.ErrorNotification, "Oops something went wrong, try again later").Render(r.Context(), w)
		if err != nil {
			h.slog.Error("Failed to Toaster: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := jwt.NewJwtService().Init(userRes.ID, userRes.Role.String())
	if err != nil {
		h.slog.Info("failed creating new jwt token for new user")
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
	w.WriteHeader(http.StatusCreated)
}
