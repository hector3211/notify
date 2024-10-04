package handlers

import (
	"log"
	"net/http"
	"server/views/templates"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var isAuth templates.IsAuthenticated = false
	page := templates.Home()
	err := templates.Layout(page, isAuth, "Notify").Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to home page: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
