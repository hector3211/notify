package handlers

import (
	"log/slog"
	"net/http"
	"server/middleware"
	"server/views/templates"
)

type LoginHandler struct {
	slog *slog.Logger
}

func NewLoginHanlder(slog *slog.Logger) *LoginHandler {
	return &LoginHandler{slog: slog}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var isAuth templates.IsAuthenticated = false
	userCtx := middleware.GetUserCtxFromCookie(w, r)
	if userCtx != nil {
		w.Header().Set("HX-Redirect", "/profile")
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
	// log.Printf("userCTX: %+v", userCtx)

	page := templates.Login()
	err := templates.Layout(page, isAuth, "login", "Notify-login").Render(r.Context(), w)
	if err != nil {
		h.slog.Error("Failed to render login page: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
