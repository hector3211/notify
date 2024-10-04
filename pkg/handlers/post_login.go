package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"server/pkg/jwt"
	"server/pkg/service"
	"server/utils"
	"strings"
	"time"
)

type PostLoginHandler struct {
	db *sql.DB
}

func NewPostLoginHandler(db *sql.DB) *PostLoginHandler {
	return &PostLoginHandler{
		db: db,
	}
}

func (h *PostLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logInEmail := r.FormValue("email")
	logInPassword := r.FormValue("password")
	userService := service.NewUserService(h.db)

	log.Printf("login Password: %s\n", logInPassword)

	if !userService.CheckEmailExists(strings.TrimSpace(logInEmail)) {
		log.Printf("No such email exists in db %s\n", logInEmail)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := userService.GetUserByEmail(strings.TrimSpace(logInEmail))
	if user == nil {
		log.Printf("No User found with email %s\n", logInEmail)
		http.Error(w, "no user found with that email", http.StatusNotFound)
		return

	}

	if !userService.CheckPasswordHash(strings.TrimSpace(logInPassword), user.ID) {
		log.Println("passwords did NOT match with hash")
		http.Error(w, "Passwords didnt match", http.StatusUnauthorized)
		return
	}

	token, err := jwt.NewJwtService().Init(user.ID, user.Role.String())
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
	w.WriteHeader(http.StatusSeeOther)
}
