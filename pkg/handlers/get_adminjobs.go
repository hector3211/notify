package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"server/middleware"
	"server/models"
	"server/pkg/service"
	"server/views/templates"
)

type AdminJobHandler struct {
	db *sql.DB
}

func NewAdminJobHandler(db *sql.DB) *AdminJobHandler {
	return &AdminJobHandler{db: db}
}

func (h *AdminJobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userCtx, ok := middleware.GetUserCtx(r)
	if !ok {
		log.Println("User context not found")
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if userCtx.ID == 0 {
		log.Printf("User context is missing ID: %+v", userCtx)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if userCtx.Role != models.ADMIN {
		log.Printf("user is not an admin: %+v", userCtx)
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		w.WriteHeader(http.StatusUnauthorized)
		return

	}

	invoiceService := service.NewInvoiceService(h.db)
	jobs := invoiceService.GetAllInvoices()

	page := templates.AdminJobs(jobs)
	err := templates.Admin(page, *userCtx, "Notify-admin").Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render admin page: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
