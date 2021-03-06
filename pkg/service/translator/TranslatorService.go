package translator

import (
	"regexp"
)

// Languages
const EnLanguage = "en"
const RuLanguage = "ru"

type TranslatorService struct {
	gateway *TranslatorGateway
}

// NewTranslatorService - constructor of TranslatorService
func NewTranslatorService(gateway *TranslatorGateway) *TranslatorService {
	return &TranslatorService{
		gateway: gateway,
	}
}

// TranslateText - translating text and auto-detecting of source and target languages
func (translator *TranslatorService) TranslateText(text string) (string, error) {
	sourceLanguage, targetLanguage := translator.detectLanguages(text)

	return translator.gateway.GetTranslation(sourceLanguage, targetLanguage, text)
}

// DetectLanguages - detecting source and target language by regex
func (translator *TranslatorService) detectLanguages(text string) (string, string) {
	if matched, _ := regexp.MatchString("[а-яА-Я]+", text); matched {
		return RuLanguage, EnLanguage
	}

	return EnLanguage, RuLanguage
}
