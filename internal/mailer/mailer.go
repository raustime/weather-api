package mailer

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendConfirmationEmail(to, token string) error {
	from := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	host := os.Getenv("SMTP_HOST") // наприклад, smtp.gmail.com
	port := os.Getenv("SMTP_PORT") // 587 або 465
	addr := fmt.Sprintf("%s:%s", host, port)

	msg := fmt.Sprintf(`Subject: Confirm your subscription

Click the link to confirm your subscription:
http://localhost:8080/api/confirm/%s
`, token)

	auth := smtp.PlainAuth("", from, pass, host)
	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}
