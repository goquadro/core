package core

import (
	"bytes"
	"html/template"
	"log"

	mailgun "github.com/mailgun/mailgun-go"
	"gopkg.in/mgo.v2/bson"
)

const (
	SubscribersCollection       = "newsletter"
	SIGNUP_NOTIFICATION_MESSAGE = "Hi, {{.Username}}! You've just signed up to GoQuadro.\nPlease click here to confirm your email address: {{.ConfirmationCode}}"
)

// ValidateEmailAddress is a wrapper for Mailgun email validator
func ValidateEmailAddress(email string) (mailgun.EmailVerification, error) {
	gun := mailgun.NewMailgun(gqConfig.mailgunDomain, gqConfig.mailgunKey, gqConfig.mailgunPubKey)
	return gun.ValidateEmail(email)
}

// Add a single email address to the mailing list.
// ToDo: Implement mailgun mailing list system
func RegisterToNewsletter(email string) error {
	locSession := getSession()
	defer locSession.Close()
	c := locSession.DB(gqConfig.jobDatabase).C(SubscribersCollection)
	err := c.Insert(bson.M{"email": email})
	return err
}

// Wrapper for the Mailgun sender API
func SendMail(subject, body, recipient string) error {
	gun := mailgun.NewMailgun(gqConfig.mailgunDomain, gqConfig.mailgunKey, gqConfig.mailgunPubKey)

	m := mailgun.NewMessage(gqConfig.notificationAddress, subject, body, recipient)
	response, id, err := gun.Send(m)
	log.Printf("Response ID: %s\n", id)
	log.Printf("Message from server: %s\n", response)
	return err
}

// Send a confirmation email to the newly registered user.
func (u User) SendConfirmationEmail() error {
	var signup_notification_message []byte
	buf := bytes.NewBuffer(signup_notification_message)
	t := template.Must(template.New("letter").Parse(SIGNUP_NOTIFICATION_MESSAGE))
	t.Execute(buf, u)
	return SendMail("You just registered on GoQuadro", string(signup_notification_message), u.Email)
}
