package email

import (
	"net/smtp"

	"github.com/tula3and/me-sign/utils"
)

func Send(to string) {
	// https://support.google.com/mail/answer/7126229?p=BadCredentials
	// Check this part before uploading to github
	auth := smtp.PlainAuth("", "<send>", "<pw>", "smtp.gmail.com")

	from := "<send>"
	tos := []string{to}

	headerSubject := "Subject: <title>\r\n"
	headerBlank := "\r\n"

	body := "<body>"

	msg := []byte(headerSubject + headerBlank + body)

	err := smtp.SendMail("smtp.gmail.com:587", auth, from, tos, msg)
	utils.HandleErr(err)
}
