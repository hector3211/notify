package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"server/middleware"
	"server/views/templates"
)

type GetAdminCreateJobHandler struct {
	db *sql.DB
}

func NewGetAdminCreateJobHandler(db *sql.DB) *GetAdminCreateJobHandler {
	return &GetAdminCreateJobHandler{db: db}
}

func (h *GetAdminCreateJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserCtx(r)
	if !ok || userCtx == nil {
		log.Println("User context not found")
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	page := templates.CreateJobForm()
	err := templates.Admin(page, *userCtx, "Notify-admin").Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render admin page: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
