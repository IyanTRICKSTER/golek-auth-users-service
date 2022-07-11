package notification

import (
	"github.com/joho/godotenv"
	"log"
	"net/smtp"
	"os"
)

func Sendmail(error chan error, message string) {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	HOST := os.Getenv("MAIL_HOST")
	PASS := os.Getenv("MAIL_PASSWORD")
	USER := os.Getenv("MAIL_USERNAME")
	PORT := os.Getenv("MAIL_PORT")

	// Choose auth method and set it up
	auth := smtp.PlainAuth("", USER, PASS, HOST)

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{"bill@gates.com"}
	msg := []byte("To: bill@gates.com\r\n" +
		"Subject: Acourse Reset Password \r\n" +
		"\r\n" +
		"Here is your reset token " + message)
	err = smtp.SendMail(HOST+":"+PORT, auth, "iyan@mail.com", to, msg)

	if err != nil {
		error <- err
	}
	error <- nil
}
