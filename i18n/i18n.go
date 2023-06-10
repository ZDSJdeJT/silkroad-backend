package i18n

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	ZhCN = "zh-CN"
	EnUS = "en-US"
)

const DefaultLanguage = ZhCN

var Languages = []string{DefaultLanguage, EnUS}

func initI18n() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	for _, lang := range Languages {
		bundle.MustLoadMessageFile("./locales/" + lang + ".json")
	}
	return bundle
}

var bundle = initI18n()

func GetLocalizedMessage(langTag string, messageId string) string {
	localizedMessage, _ := i18n.NewLocalizer(bundle, langTag).Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: messageId,
		},
	})
	return localizedMessage
}

func GetLocalizedMessageWithTemplate(langTag string, messageId string, templateData interface{}) string {
	localizedMessage, _ := i18n.NewLocalizer(bundle, langTag).Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: messageId,
		},
		TemplateData: templateData,
	})
	return localizedMessage
}
