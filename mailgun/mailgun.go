package mailgun

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type MailData struct {
	Text    string `json:"text"`
	HTML    string `json:"html"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Type    string `json:"type"`
}

type MailConfig struct {
	From     string
	Domain   string
	Username string
	Password string
}

func NewMailGun(from, domain, userName, password string) MailConfig {
	return MailConfig{
		From:     from,
		Domain:   domain,
		Username: userName,
		Password: password,
	}
}

func (conn *MailConfig) Send(to, subject, text, html string) (err error) {
	mg := mailgun.NewMailgun(conn.Domain, conn.Password)
	message := mg.NewMessage(conn.From, subject, text, to)
	if html != "" {
		message.SetHtml(html)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
	return nil
}
