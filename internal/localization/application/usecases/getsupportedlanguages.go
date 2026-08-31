package usecases

import (
	"gomander/internal/localization"
)

type GetSupportedLanguages struct {
	catalogue localization.Catalogue
}

func NewGetSupportedLanguages(catalogue localization.Catalogue) *GetSupportedLanguages {
	return &GetSupportedLanguages{
		catalogue: catalogue,
	}
}

func (uc *GetSupportedLanguages) Execute() ([]string, error) {
	return uc.catalogue.Locales()
}
