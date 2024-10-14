package handlers

import (
	"log/slog"
	"net/http"
	"server/middleware"
	"server/views/templates"
)

type HomeHandler struct{ slog *slog.Logger }

func NewHomeHandler(slog *slog.Logger) *HomeHandler {
	return &HomeHandler{slog: slog}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var isAuth templates.IsAuthenticated = false
	userCtx := middleware.GetUserCtxFromCookie(w, r)
	if userCtx != nil {
		isAuth = true
	}

	page := templates.Home()
	err := templates.Layout(page, isAuth, "Notify").Render(r.Context(), w)
	if err != nil {
		h.slog.Error("Failed to home page: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
