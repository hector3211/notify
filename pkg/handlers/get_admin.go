package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"server/middleware"
	"server/views/templates"
)

type AdminHandler struct {
	db *sql.DB
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserCtx(r)
	if !ok {
		log.Println("User context not found")
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	page := templates.AdminHome()
	err := templates.Admin(page, *userCtx, "Notify-admin").Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render admin page: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
