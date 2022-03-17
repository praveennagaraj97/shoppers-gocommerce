package models

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/i18n"
)

type Mail struct {
	Sender  string
	To      []string
	Subject string
}

func (m *Mail) GenerateHTMLMail(body *bytes.Buffer) *bytes.Buffer {
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"

	msg := []byte(fmt.Sprintf("Subject: %s\n", m.Subject) +
		fmt.Sprintf("From: %v\n", strings.Join(m.To, ",")) +
		mime + "\n")

	body.Write(msg)

	return body

}

// Email Templates

type MailMetaInfo struct {
	Title     string
	BrandInfo string
}

type NewUserTemplateData struct {
	MetaInfo     *MailMetaInfo
	WelcomeTitle string
	WelcomeDesc  string
	RedirectLink string
	ButtonName   string
}

func WelcomeEmailAndVerifyEmailTemplate(a *i18n.Internationalization, c *gin.Context, redirectLink string) interface{} {
	var td NewUserTemplateData = NewUserTemplateData{
		MetaInfo: &MailMetaInfo{
			a.GetMessage("welcome_meta_title", c),
			a.GetMessage("meta_brand_info", c),
		},
		WelcomeTitle: a.GetMessage("welcome_template_title", c),
		WelcomeDesc:  a.GetMessage("welcome_template_desc", c),
		ButtonName:   a.GetMessage("confirm_email_address", c),
		RedirectLink: redirectLink,
	}
	return td
}

type ForgotPasswordTemplateData struct {
	MetaInfo    *MailMetaInfo
	Title       string
	Description string
	Link        string
	ButtonName  string
}

func ForgotPasswordEmailTemplateData(a *i18n.Internationalization, c *gin.Context, redirectLink string) interface{} {
	td := &ForgotPasswordTemplateData{
		MetaInfo: &MailMetaInfo{
			a.GetMessage("password_reset", c),
			a.GetMessage("meta_brand_info", c),
		},
		Title:       a.GetMessage("password_reset", c),
		Description: a.GetMessage("password_reset_requested", c),
		Link:        redirectLink,
		ButtonName:  "continue",
	}

	return td
}
