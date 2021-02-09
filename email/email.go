package email

import (
	"bytes"
	"fmt"
	"github.com/arman-aminian/type-your-song/secure"
	"net/smtp"
	"text/template"
)

func SendEmail(to []string, text string) error {
	// Sender data.
	from := "typeyoursong@gmail.com"
	password := secure.EmailPassword

	// Receiver email address.
	//to := []string{
	//	"sender@example.com",
	//}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("./email/template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: This is a test subject \n%s\n\n", mimeHeaders)))

	err := t.Execute(&body, struct {
		Name    string
		Message string
	}{
		Name:    "Type A Song",
		Message: text,
	})

	if err != nil {
		panic(err)
		return err
	}

	// Sending email.
	fmt.Println("sending email to", to, "...")
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		panic(err)
		return err
	}
	fmt.Println("Email Sent!")
	return nil
}
