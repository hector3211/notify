package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"server/models"
	"server/pkg/service"
	"server/utils"
	"server/views/templates"
)

type SearchUserHandler struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewPostSearchUserHandler(db *sql.DB, slog *slog.Logger) *SearchUserHandler {
	return &SearchUserHandler{db: db}
}

func (h *SearchUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userService := service.NewUserService(h.db)
	query := r.FormValue("user-query")

	allUsers := userService.GetAllUsers()
	var filteredResults []models.User

	for _, res := range allUsers {
		if query == "" || utils.ContainsQuery(res.LastName, query) {
			filteredResults = append(filteredResults, res)
		}
	}

	w.Header().Set("Content-Type", "text/html")
	err := templates.UserSearch(filteredResults).Render(r.Context(), w)
	if err != nil {
		h.slog.Error("failed rendering admin page:" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
