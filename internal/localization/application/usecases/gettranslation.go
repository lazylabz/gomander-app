package usecases

import (
	"embed"
	"encoding/json"
	"fmt"

	"gomander/internal/localization/domain"
)

type GetTranslation interface {
	Execute(locale string) (*domain.Localization, error)
}

type DefaultGetTranslation struct {
	localeFs embed.FS
}

func NewGetTranslation(localeFs embed.FS) *DefaultGetTranslation {
	return &DefaultGetTranslation{
		localeFs: localeFs,
	}
}

func (uc *DefaultGetTranslation) Execute(locale string) (*domain.Localization, error) {
	localeJson, err := uc.localeFs.ReadFile(fmt.Sprintf("locales/%s.json", locale))
	if err != nil {
		return nil, fmt.Errorf("read locale json: %w", err)
	}

	var lng domain.Localization
	if err := json.Unmarshal(localeJson, &lng); err != nil {
		return nil, fmt.Errorf("unmarshal locale json: %w", err)
	}

	return &lng, nil
}