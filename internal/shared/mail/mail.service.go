package mail

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

type Service struct {
	config SMTPConfig
}

type MailService interface {
	SendPasswordResetEmail(
		to string,
		resetURL string,
	) error
}

func NewService(config SMTPConfig) *Service {
	return &Service{
		config: config,
	}
}

func (s *Service) SendPasswordResetEmail(
	to string,
	code string,
) error {

	m := gomail.NewMessage()

	m.SetHeader("From", s.config.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Recuperación de contraseña - KAI")

	body := fmt.Sprintf(`
		<h2>Recuperación de contraseña</h2>
		<p>Tu código de recuperación es:</p>
		<h1 style="letter-spacing: 4px;">%s</h1>
		<p>Este código expirará en 10 minutos.</p>
		<p>Si no solicitaste este cambio, podés ignorar este mensaje.</p>
	`, code)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		s.config.Host,
		s.config.Port,
		s.config.User,
		s.config.Password,
	)

	return d.DialAndSend(m)
}
