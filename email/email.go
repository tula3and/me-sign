package email

import (
	"fmt"
	"io/ioutil"
	"net/smtp"
	"strings"

	"github.com/tula3and/me-sign/sign"
	"github.com/tula3and/me-sign/utils"
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
	data, err := ioutil.ReadFile("mainEmail.txt")
	utils.HandleErr(err)

	mainEmail := strings.Split(string(data), "/")
	from := mainEmail[0]
	pw := mainEmail[1]

	auth := smtp.PlainAuth("", from, pw, "smtp.gmail.com")

	tos := []string{to}

	headerSubject := "Subject: [MeSign] Email Verification\r\n"
	headerBlank := "\r\n"
	emailHex := fmt.Sprintf("%x", to)
	body := "http://localhost:4000/key?email=" + emailHex + "&signed=" + sign.Sign(emailHex, sign.Key())

	msg := []byte(headerSubject + headerBlank + body)

	err = smtp.SendMail("smtp.gmail.com:587", auth, from, tos, msg)
	utils.HandleErr(err)
}
