package handlers

import (
	"log"
	"net/http"
	"server/middleware"
	"server/views/templates"
)

type SingupHandler struct{}

func NewSignupHandler() *SingupHandler {
	return &SingupHandler{}
}

func (s *SingupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var isAuth templates.IsAuthenticated
	_, ok := middleware.GetUserCtxFromCookie(w, r)
	if ok {
		w.Header().Set("HX-Redirect", "/profile")
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
	// log.Printf("userCTX: %+v", userCtx)
	// Data default
	isAuth = false
	page := templates.Signup()
	err := templates.Layout(page, isAuth, "Notify-signup").Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render sing up page: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
