package service

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/smtp"
	"os"
	"sync"
)

var (
	ErrNotifyHostEmailNil    = errors.New("no host email provided")
	ErrNotifyHostPasswordNil = errors.New("no host password provided")
	ErrEmailFailed           = errors.New("email failed to send")
)

type EmailService struct {
	db   *sql.DB
	slog *slog.Logger
	wG   *sync.WaitGroup
}

func NewEmailService(db *sql.DB, slog *slog.Logger) *EmailService {
	return &EmailService{
		db:   db,
		slog: slog,
		wG:   &sync.WaitGroup{},
	}
}

func (e EmailService) SendEmail(maxWorkers int, recipients []string, subject, msg string) error {
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

	workerChan := make(chan struct{}, maxWorkers)
	for _, recipient := range recipients {
		e.wG.Add(1)

		select {
		case workerChan <- struct{}{}:
			go func(recipient string) {
				defer e.wG.Done()
				defer func() { <-workerChan }()

				formattedBody := fmt.Sprintf("Hello, \n\n%s!\n\nBest regards,\nNotify", msg)
				msg := []byte("Subject: " + subject + "\r\n" + "\r\n" + formattedBody + "\r\n")
				auth := smtp.PlainAuth("", NotifyHostEmail, NotifyHostPassword, smtHost)
				err := smtp.SendMail(smtHost+":"+smtPort, auth, NotifyHostEmail, []string{recipient}, msg)
				if err != nil {
					e.slog.Error("failed to send email", "error", err, "recipient", recipient)
				} else {
					e.slog.Info("Successfully sent email")
				}
			}(recipient)
		default:
			e.slog.Warn("No workers available, retrying in a bit...", "recipient", recipient)
		}
	}
	e.wG.Wait()

	return nil
}
