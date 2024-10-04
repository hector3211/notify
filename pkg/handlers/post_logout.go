package handlers

import (
	"database/sql"
	"net/http"
	"time"
)

type PostLogOutHandler struct {
	db *sql.DB
}

func NewPostLogOutHandler() *PostLogOutHandler {
	return &PostLogOutHandler{}
}

func (h *PostLogOutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		MaxAge:  -1,
		Expires: time.Now().Add(-100 * time.Hour),
		Path:    "/",
	})
	w.Header().Del("Authorization")
	w.Header().Set("HX-Redirect", "/")
	http.Redirect(w, r, "/", http.StatusOK)
}
