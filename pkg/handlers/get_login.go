package handlers

import (
	"log"
	"net/http"
	"server/middleware"
	"server/views/templates"
)

type LoginHandler struct{}

func NewLoginHanlder() *LoginHandler {
	return &LoginHandler{}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var isAuth templates.IsAuthenticated
	_, ok := middleware.GetUserCtxFromCookie(w, r)
	if ok {
		w.Header().Set("HX-Redirect", "/profile")
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
	// log.Printf("userCTX: %+v", userCtx)

	// Admin default
	isAuth = false
	page := templates.Login()
	err := templates.Layout(page, isAuth, "Notify-login").Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render login page: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
