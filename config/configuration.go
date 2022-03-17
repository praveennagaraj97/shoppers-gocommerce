package conf

import (
	awspkg "github.com/praveennagaraj97/shoppers-gocommerce/pkg/aws"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/i18n"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/mail"

	"go.mongodb.org/mongo-driver/mongo"
)

type GlobalConfiguration struct {
	Port                  string
	Env                   string
	DatabaseConnectionURL string
	DatabaseName          string
	Domain                string
	DatabaseClient        *mongo.Client
	Mailer                *mail.Mailer
	Localize              *i18n.Internationalization
	FrontendBaseUrl       string
	AWSUtils              *awspkg.AWSConfiguration
	Locales               *[]string
	DefaultLocale         string
}
