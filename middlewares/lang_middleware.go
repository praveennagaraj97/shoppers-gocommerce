package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/praveennagaraj97/shoppers-gocommerce/constants"
)

// accepts lang via query param | Accept-Language Header.
func (m *Middlewares) AcceptLanguageParser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var locale string = c.Request.URL.Query().Get("locale")

		if locale == "" {
			locale = c.Request.Header.Get("Accept-Language")
			if locale == "" {
				locale = m.conf.DefaultLocale
			}
		}

		c.Set(string(constants.CUSTOME_HEADER_LANG_KEY), locale)
		c.Next()
	}
}
