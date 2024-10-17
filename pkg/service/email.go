package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/smtp"
	"os"
)

var (
	ErrNotifyHostEmailNil    = errors.New("no host email provided")
	ErrNotifyHostPasswordNil = errors.New("no host password provided")
	ErrEmailFailed           = errors.New("email failed to send")
)

type EmailService struct {
	db   *sql.DB
	slog *slog.Logger
}

func NewEmailService(db *sql.DB, slog *slog.Logger) *EmailService {
	return &EmailService{
		db:   db,
		slog: slog,
	}
}

func (e EmailService) SendEmail(maxWorkers int, recipientsEmail, subject, msg string) error {
	NotifyHostEmail := os.Getenv("NOTIFY_HOST_EMAIL")
	NotifyHostPassword := os.Getenv("NOTIFY_HOST_PASSWORD")
	if NotifyHostEmail == "" {
		return ErrNotifyHostEmailNil
	}
	if NotifyHostPassword == "" {
		return ErrNotifyHostPasswordNil
	}

	smtHost := "smtp.gmail.com"
	smtPort := "587"

	formattedBody := fmt.Sprintf("Hello, \n\n%s!\n\nBest regards,\nNotify", msg)
	email := []byte("Subject: " + subject + "\r\n" + "\r\n" + formattedBody + "\r\n")
	auth := smtp.PlainAuth("", NotifyHostEmail, NotifyHostPassword, smtHost)
	err := smtp.SendMail(smtHost+":"+smtPort, auth, NotifyHostEmail, []string{recipientsEmail}, email)
	if err != nil {
		e.slog.Error("failed to send email", "error", err, "recipient", recipientsEmail)
	} else {
		e.slog.Info("Successfully sent email")
	}

	return nil
}
