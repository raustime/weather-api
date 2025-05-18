package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"weatherapi/internal/openweatherapi"
)

type confirmationData struct {
	ConfirmURL string
}

func SendConfirmationEmail(to, token string) error {

	apiHost := os.Getenv("APP_BASE_URL")
	link := fmt.Sprintf("%s/api/confirm/%s", apiHost, token)
	data := confirmationData{ConfirmURL: link}

	tmpl, err := template.ParseFiles("internal/templates/confirmation_email.html")
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}
	return sendHTMLEmail(to, "Confirm your subscription", body.String())
}

func SendWeatherEmail(to string, city string, weather *openweatherapi.WeatherData, baseURL, token string) error {
	tmpl, err := template.ParseFiles("internal/templates/weather_email.html")
	if err != nil {
		return err
	}

	var body bytes.Buffer
	data := struct {
		City           string
		Description    string
		Temperature    float64
		Humidity       float64
		UnsubscribeURL string
	}{
		City:           city,
		Description:    weather.Description,
		Temperature:    weather.Temperature,
		Humidity:       weather.Humidity,
		UnsubscribeURL: fmt.Sprintf("%s/api/unsubscribe/%s", baseURL, token),
	}
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	subject := fmt.Sprintf("Weather Update for %s", city)
	return sendHTMLEmail(to, subject, body.String())
}
func sendHTMLEmail(to, subject, htmlBody string) error {
	from := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	msg := fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
		htmlBody

	auth := smtp.PlainAuth("", from, pass, host)
	err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("Error sending HTML email to %s: %v", to, err)
	}
	return err
}
