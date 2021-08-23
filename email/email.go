package email

import (
	"net/smtp"
	"strings"

	"github.com/tula3and/me-sign/utils"
)

const (
	FROM string = ""
)

func Verify(address string) bool {
	div := strings.Split(address, "@")
	if len(div) != 2 {
		return false
	}
	web := strings.Split(div[1], ".")
	if len(web) != 2 {
		return false
	}
	send(address)
	return true
}

func send(to string) {
	// https://support.google.com/mail/answer/7126229?p=BadCredentials
	// Check this part before uploading to github
	auth := smtp.PlainAuth("", FROM, "<pw>", "smtp.gmail.com")

	tos := []string{to}

	headerSubject := "Subject: [MeSign] Email Verification\r\n"
	headerBlank := "\r\n"
	body := "<body>"

	msg := []byte(headerSubject + headerBlank + body)

	err := smtp.SendMail("smtp.gmail.com:587", auth, FROM, tos, msg)
	utils.HandleErr(err)
}
