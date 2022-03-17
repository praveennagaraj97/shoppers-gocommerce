package cmd

import (
	conf "github.com/praveennagaraj97/shoppers-gocommerce/config"
	"github.com/praveennagaraj97/shoppers-gocommerce/db"
	awspkg "github.com/praveennagaraj97/shoppers-gocommerce/pkg/aws"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/env"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/i18n"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/mail"
)

// initializeApp provides Default / initial configuration for the application.
func initializeApp() *conf.GlobalConfiguration {

	var mongo_url = env.GetEnvVariable("MONGO_URI")

	// translation package
	locales := []string{"en", "fr", "ar"}
	localize := &i18n.Internationalization{}
	localize.Initialize(locales)

	// mailer package
	m := &mail.Mailer{}
	m.Initialize()

	awsPkg := &awspkg.AWSConfiguration{}
	awsPkg.Initialize()

	app := &conf.GlobalConfiguration{
		Port:                  env.GetEnvVariable("PORT"),
		Env:                   env.GetEnvVariable("ENVIRONMENT"),
		FrontendBaseUrl:       "http://localhost:8080",
		DatabaseConnectionURL: mongo_url,
		DatabaseName:          "shoppers-gocommerce",
		Domain:                "localhost",
		Localize:              localize,
		Mailer:                m,
		AWSUtils:              awsPkg,
		Locales:               &locales,
		DefaultLocale:         "en",
	}

	// DB client stores the ref of database client.
	app.DatabaseClient = db.InitializeDatabaseConnection(app.DatabaseConnectionURL, app.DatabaseName)

	return app
}
