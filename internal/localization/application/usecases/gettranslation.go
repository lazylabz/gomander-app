package usecases

import (
	"gomander/internal/localization"
)

type GetTranslation struct {
	catalogue localization.Catalogue
}

func NewGetTranslation(catalogue localization.Catalogue) *GetTranslation {
	return &GetTranslation{
		catalogue: catalogue,
	}
}

func (uc *GetTranslation) Execute(locale string) (*localization.Localization, error) {
	return uc.catalogue.Translations(locale)
}
