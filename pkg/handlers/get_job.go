package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"server/pkg/service"
	"server/views/templates"
	"strconv"
	"strings"
)

type JobHander struct {
	db *sql.DB
}

func NewGetJobHandler(db *sql.DB) *JobHander {
	return &JobHander{db: db}
}

func (h *JobHander) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userEmail := r.FormValue("email")
	userInvoice := r.FormValue("invoice")
	userService := service.NewUserService(h.db)
	invoiceService := service.NewInvoiceService(h.db)

	if !userService.CheckEmailExists(strings.TrimSpace(userEmail)) {
		fmt.Printf("No such email exists in db %s", userEmail)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := userService.GetUserByEmail(userEmail)
	if user == nil {
		fmt.Printf("No such user exists in db %s", userEmail)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	InvoiceIdNum, err := strconv.Atoi(userInvoice)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	invoiceId := invoiceService.GetInvoiceId(user.ID, InvoiceIdNum)
	if invoiceId == 0 {
		fmt.Printf("No such invoice exists in db %s", userInvoice)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: fix string to int here
	invoice := invoiceService.GetInvoice(string(invoiceId))
	if invoice == nil {
		fmt.Printf("No such job exists in db with invoice %s", userInvoice)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = templates.Invoice(*invoice).Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to job section: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
