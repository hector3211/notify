package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/models"
	"server/pkg/jwt"
	"server/pkg/service"
	"server/utils"
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

	h.slog.Info("signup email: %s\n", signupEmail)
	h.slog.Info("signup password: %s\n", signupPassword)

	userService := service.NewUserService(h.db)
	if userService.CheckEmailExists(strings.TrimSpace(signupEmail)) {
		http.Error(w, "user with email already exists", http.StatusConflict)
		return
	}

	user := &models.User{
		FirstName: signupFirstName,
		LastName:  signupLastName,
		Password:  signupPassword,
		Email:     signupEmail,
		Role:      models.UserRole(models.USER),
	}

	// NOTE: CreateUser func takes care of hashing the password
	userRes := userService.CreateUser(user)
	if userRes == nil {
		http.Error(w, "failed creating new user", http.StatusInternalServerError)
		return
	}

	token, err := jwt.NewJwtService().Init(userRes.ID, userRes.Role.String())
	if err != nil {
		http.Error(w, "failed creating jwt token", http.StatusInternalServerError)
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

	// w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w.Header().Set("HX-Redirect", "/profile")
	w.WriteHeader(http.StatusCreated)
}
