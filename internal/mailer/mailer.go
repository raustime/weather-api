package mailer

import (
	"bytes"
	"fmt"
	"html/template"

	"os"
	"weatherapi/internal/openweatherapi"
)

type confirmationData struct {
	ConfirmURL string
}

var Email EmailSender = &SMTPSender{} // дефолтно справжній

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
	return Email.Send(to, "Confirm your subscription", body.String())
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
	return Email.Send(to, subject, body.String())
}
