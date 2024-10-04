package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"server/models"
	"server/pkg/service"
	"server/utils"
	"server/views/templates"
)

type SearchUserHandler struct {
	db *sql.DB
}

func NewPostSearchUserHandler(db *sql.DB) *SearchUserHandler {
	return &SearchUserHandler{db: db}
}

func (h *SearchUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	invoiceService := service.NewUserService(h.db)
	query := r.FormValue("user-query")

	allUsers := invoiceService.GetAllUsers()
	var filteredResults []models.User

	for _, res := range allUsers {
		if query == "" || utils.ContainsQuery(res.LastName, query) {
			filteredResults = append(filteredResults, res)
		}
	}

	w.Header().Set("Content-Type", "text/html")
	err := templates.UserSearch(filteredResults).Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render admin page: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
