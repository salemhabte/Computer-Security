package infrastructure

import (
	"fmt"
	"math/rand"
	"net/smtp"
)

type EmailService struct {
	From        string
	AppPass     string
	SmtpService string
	SmtPort     string
	Username    string
}

func NewOTP_service(from, app, smtpsev, smtport, username string) *EmailService {
	return &EmailService{
		From:        from,
		AppPass:     app,
		SmtpService: smtpsev,
		SmtPort:     smtport,
		Username:    username,
	}
}

func (s *EmailService) GenerateRandomOTP() string {
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))
	return otp
}
func (s *EmailService) Send(toEmail, otp string) error {

	auth := smtp.PlainAuth("", s.From, s.AppPass, s.SmtpService)
	msg := []byte(fmt.Sprintf("Subject: OTP Verification\n\nYour OTP is: %s", otp))
	return smtp.SendMail(s.SmtpService+":"+s.SmtPort, auth, s.From, []string{toEmail}, msg)
}

func (s *EmailService) SendResetLink(toEmail, subject, message string) error {

	auth := smtp.PlainAuth("", s.From, s.AppPass, s.SmtpService)
	msg := []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, message))
	return smtp.SendMail(s.SmtpService+":"+s.SmtPort, auth, s.From, []string{toEmail}, msg)
}

