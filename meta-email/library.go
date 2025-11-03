package metaemail

import (
	"fmt"
	"net/smtp"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func IsEmailValid(email string) bool {
	return emailRegex.MatchString(email)
}

func MaskEmail(s string) string {
	n := len(s)
	switch {
	case n <= 2:
		return s
	case n == 3:
		return s[:1] + "*" + s[2:]
	case n == 4:
		return s[:1] + "**" + s[3:]
	default:
		// 中间星号数量为 min(n-4, 4)
		numStars := n - 4
		if numStars > 4 {
			numStars = 4
		}
		return s[:2] + strings.Repeat("*", numStars) + s[n-2:]
	}
}

func SendEmail(
	nickname string,
	userEmail string,
	userPassword string,
	host string,
	port int,
	to string,
	subject string,
	body string,
) error {
	address := fmt.Sprintf("%s:%d", host, port)
	msg := []byte("From: " + nickname + "<" + userEmail + ">\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" + body)
	auth := smtp.PlainAuth("", userEmail, userPassword, host)
	return smtp.SendMail(address, auth, userEmail, []string{to}, msg)
}
