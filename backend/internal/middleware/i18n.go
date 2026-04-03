package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	DefaultLanguage = "en"
	LanguageContextKey = "lang"
	LocalizerContextKey = "localizer"
)

// I18nMiddleware detects the language from the Accept-Language header
// and sets the localizer in the Gin context.
func I18n(bundle *i18n.Bundle) gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptLang := c.GetHeader("Accept-Language")
		lang := c.DefaultQuery("lang", acceptLang)

		// Try to match the language tag
		tags, _, err := language.ParseAcceptLanguage(lang)
		if err != nil || len(tags) == 0 {
			tags = []language.Tag{language.Make(DefaultLanguage)}
		}

		matcher := language.NewMatcher([]language.Tag{
			language.Make("en"),
			language.Make("zh"),
		})
		_, idx, _ := matcher.Match(tags...)

		var matchedLang string
		switch idx {
		case 0:
			matchedLang = "en"
		case 1:
			matchedLang = "zh"
		default:
			matchedLang = DefaultLanguage
		}

		localizer := i18n.NewLocalizer(bundle, matchedLang)

		c.Set(LanguageContextKey, matchedLang)
		c.Set(LocalizerContextKey, localizer)
		c.Next()
	}
}

// GetLocalizer retrieves the localizer from the Gin context.
func GetLocalizer(c *gin.Context) *i18n.Localizer {
	if val, exists := c.Get(LocalizerContextKey); exists {
		if l, ok := val.(*i18n.Localizer); ok {
			return l
		}
	}
	return nil
}

// Localize is a helper to localize a message.
func Localize(c *gin.Context, messageID string) string {
	localizer := GetLocalizer(c)
	if localizer == nil {
		return messageID
	}
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		return messageID
	}
	return msg
}
