package mail

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/color"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/env"
	logger "github.com/praveennagaraj97/shoppers-gocommerce/pkg/log"
)

type Mailer struct {
	sender struct {
		username string
		password string
		email    string
	}
	smtp struct {
		host string
		port string
	}
	address       string
	templateCache map[string]*template.Template
}

// initialize mail package and returns mail instance which can be used to send emails.
func (m *Mailer) Initialize() {

	m.sender.email = env.GetEnvVariable("SENDER_EMAIL")
	m.sender.password = env.GetEnvVariable("SMTP_PASSWORD")
	m.sender.username = env.GetEnvVariable("SMTP_USERNAME")
	m.smtp.host = env.GetEnvVariable("SMTP_HOST")
	m.smtp.port = env.GetEnvVariable("SMTP_PORT")
	m.templateCache = make(map[string]*template.Template)

	m.address = fmt.Sprintf("%s:%s", m.smtp.host, m.smtp.port)

	logger.PrintLog("Mailer Package initialized", color.White)

}

// send mail template with formated data to client
func (m *Mailer) SendNoReplyMail(to []string, subject string, templateName string, td interface{}) error {

	smtpAuth := smtp.PlainAuth("", m.sender.username, m.sender.password, m.smtp.host)

	t := m.parseTemplate(templateName)

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("To: %v \nFrom: %s \nSubject: %s \n%s\n\n", to, "noreply@shopee.com", subject, mimeHeaders)))

	err := t.Execute(&body, td)

	if err != nil {
		return err
	}

	err = smtp.SendMail(m.address, smtpAuth, m.sender.email, to, body.Bytes())

	if err != nil {
		return err
	}

	return nil
}

//go:embed templates
var emailTemplatesFS embed.FS

var funcs = template.FuncMap{}

func (m *Mailer) parseTemplate(file string) *template.Template {

	if m.templateCache[file] != nil {
		return m.templateCache[file]
	}
	t, err := template.New(
		fmt.Sprintf("%s.gotmpl", file)).Funcs(funcs).ParseFS(
		emailTemplatesFS,
		fmt.Sprintf("templates/%s.gotmpl", file))

	m.templateCache[file] = t

	if err != nil {

	}

	return t

}
