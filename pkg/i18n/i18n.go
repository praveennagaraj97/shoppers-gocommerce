package i18n

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/shoppers-gocommerce/constants"
	"github.com/praveennagaraj97/shoppers-gocommerce/pkg/color"
	logger "github.com/praveennagaraj97/shoppers-gocommerce/pkg/log"
)

type Internationalization struct {
	allowedLocales       []string
	translationsMessages map[string]map[string]string
}

//go:embed translations
var translationsFiles embed.FS

/*
Initialized translations package for application.
	locales Use ISO2 format with lower cases https://www.fao.org/countryprofiles/iso3list/en/
*/
func (i *Internationalization) Initialize(locales []string) {

	i.allowedLocales = locales
	i.translationsMessages = make(map[string]map[string]string)

	// show warning for missing translations.

	for _, lang := range i.allowedLocales {
		bytes, err := translationsFiles.ReadFile(fmt.Sprintf("translations/messages/%s.json", lang))
		if err != nil {
			logger.PrintLog(fmt.Sprintf("Translation Message file not found for given locale '%s'.", lang), color.Red)
			continue
		}
		// format to map data.
		var msgs map[string]string
		err = json.Unmarshal(bytes, &msgs)
		if err != nil {
			logger.ErrorLogFatal(err)
			return
		}
		i.translationsMessages[lang] = msgs
	}

	logger.PrintLog("Localization Package Initialized 🌐", color.Yellow)
}

// provide key of translation.
func (i *Internationalization) GetMessage(key string, c *gin.Context) string {

	locale := i.ParseLanguageFromSetHeader(c)

	// locale not found fallback to en.
	isLocaleFound := find(i.allowedLocales, locale)
	if !isLocaleFound {
		locale = "en"
	}

	messages := i.translationsMessages[locale]

	message := messages[key]

	if message == "" {
		return key
	}

	return message
}

func (i *Internationalization) ParseLanguageFromSetHeader(c *gin.Context) string {
	lang, exists := c.Get(string(constants.CUSTOME_HEADER_LANG_KEY))

	if !exists {
		return "en"
	}

	return fmt.Sprintf("%v", lang)
}

func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
