package service

import (
	"database/sql"
	"fmt"
	"log/slog"
	"server/models"
	"strconv"

	"github.com/hector3211/shogun"
)

type InvoiceService struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewInvoiceService(db *sql.DB, slog *slog.Logger) *InvoiceService {
	return &InvoiceService{db: db, slog: slog}
}

func (i InvoiceService) UpdateInvoiceStatus(newStatus models.JobStatus, invoiceId string) error {
	invID, err := strconv.Atoi(invoiceId)
	if err != nil {
		return err
	}

	query := shogun.Update("invoices").
		Set(
			shogun.Equal("status", newStatus.String()),
		).
		Where(shogun.Equal("id", invID))

	_, err = i.db.Exec(query.Build())
	if err != nil {
		return err
	}

	return nil
}

func (i InvoiceService) CreateInvoice(userId int, invoiceNum, installDate string) error {
	query := shogun.Insert("invoices").
		Columns("user_id", "invoice", "status", "install_date").
		Values(
			userId,
			invoiceNum,
			models.JOBPENDING.String(),
			installDate,
		)

	_, err := i.db.Exec(query.Build())
	if err != nil {
		return err
	}

	return nil
}

func (i InvoiceService) DeleteInvoice(invoiceIDStr string) error {
	invID, err := strconv.Atoi(invoiceIDStr)
	if err != nil {
		return err
	}

	query := shogun.Delete("invoices").
		Where(shogun.Equal("id", invID))
	// query := "DELETE FROM invoices WHERE id = ?"

	// fmt.Println(query.String())

	_, err = i.db.Exec(query.Build())
	if err != nil {
		return err
	}

	return nil
}

func (i InvoiceService) GetInvoiceId(userId int, invoiceNum int) int {
	var invoiceId int

	query := shogun.NewSelectBuilder().
		Select("id").From("invoices").
		Where(
			shogun.Equal("user_id", userId),
			shogun.Equal("invoice", invoiceNum),
		)
	err := i.db.QueryRow(query.Build()).Scan(&invoiceId)
	if err != nil {
		return 0
	}

	return invoiceId
}

func (i InvoiceService) GetAllInvoices() []models.Invoice {
	var invoices []models.Invoice

	query := shogun.NewSelectBuilder().
		Select("id", "user_id", "invoice", "status", "install_date", "created_at").
		From("invoices").
		OrderBy("created_at").Desc()

	rows, err := i.db.Query(query.Build())
	if err != nil {
		fmt.Printf("failed query invoices: %v", err)
		return invoices
	}
	defer rows.Close()

	for rows.Next() {
		var invoice models.Invoice
		err := rows.Scan(&invoice.ID, &invoice.UserId, &invoice.Invoice, &invoice.Status, &invoice.InstallDate, &invoice.CreatedAt)
		if err != nil {
			fmt.Printf("failed scanning invoices: %v", err)
		} else {
			invoices = append(invoices, invoice)
		}
	}

	return invoices
}

func (i InvoiceService) GetInvoice(invoiceId string) *models.Invoice {
	var invoice models.Invoice
	invID, err := strconv.Atoi(invoiceId)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	query := shogun.NewSelectBuilder().
		Select("id", "user_id", "invoice", "status", "created_at").
		From("invoices").
		Where(shogun.Equal("id", invID))

	err = i.db.QueryRow(query.Build()).Scan(&invoice.ID, &invoice.UserId, &invoice.Invoice, &invoice.Status, &invoice.CreatedAt)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	return &invoice
}

func (i InvoiceService) GetUserInvoices(userId int) []models.Invoice {
	var invoices []models.Invoice

	query := shogun.
		Select("id", "user_id", "invoice", "status", "created_at").
		From("invoices").
		Where(shogun.Equal("user_id", userId)).
		OrderBy("created_at").Desc()

	rows, err := i.db.Query(query.Build())
	if err != nil {
		return invoices
	}
	defer rows.Close()

	for rows.Next() {
		var inv models.Invoice

		err := rows.Scan(&inv.ID, &inv.UserId, &inv.Invoice, &inv.Status, &inv.CreatedAt)
		if err != nil {
			fmt.Printf("failed scanning for invoice: %v", err)
		} else {
			invoices = append(invoices, inv)
		}

	}

	return invoices
}
