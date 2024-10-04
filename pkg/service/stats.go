package service

import (
	"database/sql"
	"server/models"
	"time"
)

type StatsService struct {
	db *sql.DB
}

func NewStatsService(db *sql.DB) *StatsService {
	return &StatsService{
		db: db,
	}
}

func (s StatsService) GetRecentUsers() ([]models.User, error) {
	var recentUsers []models.User
	today := time.Now().Format("2006-01-02")

	query := `SELECT id, user_id, invoice, status, created_at FROM invoices WHERE DATE(created_at) = ?`

	rows, err := s.db.Query(query, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var email string
		var lastName string
		var createdAt time.Time

		if err := rows.Scan(&id, &lastName, &email, &createdAt); err != nil {
			return nil, err
		}

		if createdAt.Format("2006-01-02") == today {
			recentUsers = append(recentUsers, models.User{ID: id, LastName: lastName, Email: email, CreatedAt: createdAt})
		}
	}

	return recentUsers, nil
}

func (s StatsService) GetRecentInvoices() ([]models.Invoice, error) {
	var recentInvoices []models.Invoice
	today := time.Now().Format("2006-01-02")

	query := `SELECT id, user_id, invoice, status, created_at FROM invoices WHERE DATE(created_at) = ?`

	rows, err := s.db.Query(query, today)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var userId int
		var invoiceNumber string
		var status string
		var createdAt time.Time

		if err := rows.Scan(&id, &userId, &invoiceNumber, &status, &createdAt); err != nil {
			return nil, err
		}

		if createdAt.Format("2006-01-02") == today {
			recentInvoices = append(recentInvoices, models.Invoice{ID: id, UserId: userId, Invoice: invoiceNumber, Status: models.JobStatus(status), CreatedAt: createdAt})
		}
	}

	return recentInvoices, nil
}
