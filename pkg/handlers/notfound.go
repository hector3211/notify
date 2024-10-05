package handlers

import (
	"net/http"
	"server/middleware"
	"server/views/templates"
)

type NotFoundHandler struct{}

func NewNotFoundHanlder() *NotFoundHandler {
	return &NotFoundHandler{}
}

func (n *NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var isAuth templates.IsAuthenticated
	userCtx := middleware.GetUserCtxFromCookie(w, r)
	if userCtx == nil {
		isAuth = false
	} else {
		isAuth = true
	}
	page := templates.NotFound()
	err := templates.Layout(page, isAuth, "Notify").Render(r.Context(), w)
	if err != nil {
		return
	}
}
