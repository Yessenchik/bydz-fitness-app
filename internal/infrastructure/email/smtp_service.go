package email

import (
	"context"
	"fmt"
	"net/smtp"
	"time"
)

type Service interface {
	SendPurchaseConfirmation(ctx context.Context, to string, planName string, endDate time.Time) error
	SendExpirationWarning(ctx context.Context, to string, endDate time.Time) error
	SendExpired(ctx context.Context, to string) error
}

type SMTPService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSMTPService(
	host string,
	port string,
	username string,
	password string,
	from string,
) *SMTPService {
	return &SMTPService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (s *SMTPService) SendPurchaseConfirmation(
	ctx context.Context,
	to string,
	planName string,
	endDate time.Time,
) error {
	subject := "Gym membership purchase confirmation"

	body := fmt.Sprintf(
		"Your membership plan %s is active until %s.",
		planName,
		endDate.Format("2006-01-02"),
	)

	return s.send(ctx, to, subject, body)
}

func (s *SMTPService) SendExpirationWarning(
	ctx context.Context,
	to string,
	endDate time.Time,
) error {
	subject := "Your gym membership is expiring soon"

	body := fmt.Sprintf(
		"Your membership will expire on %s. Please extend your subscription to continue using the gym.",
		endDate.Format("2006-01-02"),
	)

	return s.send(ctx, to, subject, body)
}

func (s *SMTPService) SendExpired(
	ctx context.Context,
	to string,
) error {
	subject := "Gym membership expired"

	body := "Your gym membership has expired. Please buy or extend your subscription to get access again."

	return s.send(ctx, to, subject, body)
}

func (s *SMTPService) send(
	ctx context.Context,
	to string,
	subject string,
	body string,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	auth := smtp.PlainAuth(
		"",
		s.username,
		s.password,
		s.host,
	)

	message := []byte(
		"From: " + s.from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	address := s.host + ":" + s.port

	return smtp.SendMail(
		address,
		auth,
		s.from,
		[]string{to},
		message,
	)
}
