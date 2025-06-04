package utils

import (
	"os"

	"github.com/loid-lab/e-commerce-api/models"
	gomail "gopkg.in/mail.v2"
)

func SendMail(data models.EmailData) error {

	message := gomail.NewMessage()
	message.SetHeader("From", os.Getenv("MAIL_FROM"))
	message.SetHeader("To", data.To)
	message.SetHeader("Subject", data.Subject)
	message.SetBody("text/html", data.HTMLBody)

	if data.ImagePath != "" {
		message.Embed(data.ImagePath)
	}

	dialer := gomail.NewDialer(
		os.Getenv("MAIL_HOST"),
		587,
		os.Getenv("MAIL_USER"),
		os.Getenv("MAIL_PASS"),
	)

	return dialer.DialAndSend(message)
}
