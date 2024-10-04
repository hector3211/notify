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

type ProfileHandler struct {
	db *sql.DB
}

func NewProfileHandler(db *sql.DB) *ProfileHandler {
	return &ProfileHandler{db: db}
}

func (h *ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var isAuth templates.IsAuthenticated
	userCtx, ok := middleware.GetUserCtx(r)
	if !ok || userCtx == nil {
		w.Header().Set("HX-Redirect", "/")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if userCtx.Role == models.ADMIN {
		w.Header().Set("HX-Redirect", "/admin")
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	invoiceService := service.NewInvoiceService(h.db)
	jobs := invoiceService.GetUserInvoices(userCtx.ID)

	isAuth = true
	page := templates.Profile(jobs)
	err := templates.Layout(page, isAuth, "Notify-profile").Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render profile page: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
