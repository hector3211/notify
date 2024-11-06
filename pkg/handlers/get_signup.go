package handlers

import (
	"log/slog"
	"net/http"
	"server/middleware"
	"server/views/templates"
)

type SingupHandler struct{ slog *slog.Logger }

func NewSignupHandler(slog *slog.Logger) *SingupHandler {
	return &SingupHandler{slog: slog}
}

func (h *SingupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var isAuth templates.IsAuthenticated = false
	userCtx := middleware.GetUserCtxFromCookie(w, r)
	if userCtx != nil {
		w.Header().Set("HX-Redirect", "/profile")
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}
	// log.Printf("userCTX: %+v", userCtx)

	page := templates.Signup()
	err := templates.Layout(page, isAuth, "signup", "Notify-signup").Render(r.Context(), w)
	if err != nil {
		h.slog.Error("Failed to render sing up page: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
