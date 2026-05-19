package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"

	"go.uber.org/zap"

	"github.com/Yessenchik/bydz-fitness-app/user-auth-service/internal/domain"
)

type smtpService struct {
	host     string // smtp.gmail.com
	port     string // 587
	username string // from address
	password string // app password
	baseURL  string // https://gymapp.io — для построения ссылок
	logger   *zap.Logger
}

func NewSMTPService(host, port, username, password, baseURL string, logger *zap.Logger) domain.EmailService {
	return &smtpService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		baseURL:  baseURL,
		logger:   logger,
	}
}

func (s *smtpService) SendVerificationEmail(ctx context.Context, toEmail, firstName, token string) error {
	link := fmt.Sprintf("%s/auth/verify-email?token=%s", s.baseURL, token)

	const tmpl = `
<html><body>
<h2>Привет, {{.FirstName}}!</h2>
<p>Подтвердите ваш email, нажав на ссылку:</p>
<a href="{{.Link}}">Подтвердить email</a>
<p>Ссылка действительна 24 часа.</p>
</body></html>`

	body, err := renderTemplate(tmpl, map[string]string{
		"FirstName": firstName,
		"Link":      link,
	})
	if err != nil {
		return fmt.Errorf("verify email template: %w", err)
	}

	return s.send(toEmail, "Подтвердите ваш email — Gym App", body)
}

func (s *smtpService) SendPasswordResetEmail(ctx context.Context, toEmail, firstName, token string) error {
	link := fmt.Sprintf("%s/auth/reset-password?token=%s", s.baseURL, token)

	const tmpl = `
<html><body>
<h2>Сброс пароля, {{.FirstName}}</h2>
<p>Для сброса пароля перейдите по ссылке:</p>
<a href="{{.Link}}">Сбросить пароль</a>
<p>Ссылка действительна 1 час. Если вы не запрашивали сброс — проигнорируйте письмо.</p>
</body></html>`

	body, err := renderTemplate(tmpl, map[string]string{
		"FirstName": firstName,
		"Link":      link,
	})
	if err != nil {
		return fmt.Errorf("reset password template: %w", err)
	}

	return s.send(toEmail, "Сброс пароля — Gym App", body)
}

func (s *smtpService) send(to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	msg := []byte(fmt.Sprintf(
		"From: Gym App <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		s.username, to, subject, htmlBody,
	))

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	if err := smtp.SendMail(addr, auth, s.username, []string{to}, msg); err != nil {
		s.logger.Error("smtp send failed",
			zap.String("to", to),
			zap.String("subject", subject),
			zap.Error(err),
		)
		return fmt.Errorf("smtp send: %w", err)
	}

	s.logger.Info("email sent", zap.String("to", to), zap.String("subject", subject))
	return nil
}

func renderTemplate(tmpl string, data map[string]string) (string, error) {
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
