package exceptionist

import (
	"fmt"
	"github.com/magiconair/properties"
	"strconv"
	"strings"
)

type Language string

const (
	TR Language = "tr"
	EN Language = "en"
)

type translation struct {
	errorCode    int
	errorMessage string
}

var defaultTranslation = translation{
	errorCode:    10001,
	errorMessage: "Default.",
}

var t1 = translation{
	errorCode:    1,
	errorMessage: "Geçersiz değer.",
}

var t2 = translation{
	errorCode:    2,
	errorMessage: "Hata.",
}

type TranslationService struct {
	config Config
	translations *map[Language]bucket
}

func NewTranslationService(config Config) TranslationService {
	translations := make(map[Language]bucket)

	config = config.ensure()

	return TranslationService{
		config: config,
		translations: &translations,
	}
}

func (translationService TranslationService) AddLanguageSupport(lang Language) TranslationService {
	translations := *translationService.translations
	if _, ok := translations[lang]; !ok {
		filepath := *translationService.config.Dir + "/messages_tr.properties"
		bucket := readTranslations(filepath)
		translations[lang] = bucket
	}
	return translationService
}

func (translationService TranslationService) Translate(err ObservedError, lang Language) TranslatedError {
	translations := *translationService.translations
	bucket := translations[lang]

	if translation, ok := bucket[err.Key]; ok {
		return newTranslatedError(translation.errorCode, translation.errorMessage)
	}

	return newTranslatedError(defaultTranslation.errorCode, defaultTranslation.errorMessage)
}

func readTranslations(filepath string) bucket {
	//os.Getenv("GOPATH")+"/src/mytest/messages/"+"messages_tr.properties"
	properties := properties.MustLoadFile(filepath, properties.UTF8)

	var bucket bucket = map[string]translation{}

	for _, key := range properties.Keys() {
		val := properties.MustGet(key)
		if semiColon := strings.Index(val, ";"); semiColon >= 0 {
			errorCode, err := strconv.Atoi(val[:semiColon])
			if err != nil {
				fmt.Println("invalid errorCode in the properties file:", filepath)
			}

			errorMessage := val[semiColon+1:]
			bucket[key] = translation{
				errorCode:    errorCode,
				errorMessage: errorMessage,
			}
		}
	}

	return bucket
}
