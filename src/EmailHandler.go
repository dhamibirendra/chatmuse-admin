package src

import "gopkg.in/gomail.v2"

var (
	AlertEmail = "sajat.shrestha@gmail.com"
	fromEmail  = "chatmuse2018@gmail.com"
)

func SendEmail(subject string, message string) {
	m := gomail.NewMessage()

	m.SetHeader("From", fromEmail)
	m.SetHeader("To", AlertEmail)

	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	d := gomail.NewDialer("smtp.gmail.com", 587, fromEmail, "sajat@123")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
