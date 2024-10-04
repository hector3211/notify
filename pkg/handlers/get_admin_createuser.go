package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"server/middleware"
	"server/models"
	"server/views/templates"
)

type GetAdminCreateUserHandler struct {
	db *sql.DB
}

func NewGetAdminCreateUserHandler(db *sql.DB) *GetAdminCreateUserHandler {
	return &GetAdminCreateUserHandler{db: db}
}

func (h *GetAdminCreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	page := templates.CreateUserForm()
	err := templates.Admin(page, *userCtx, "Notify-admin").Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render admin page: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
