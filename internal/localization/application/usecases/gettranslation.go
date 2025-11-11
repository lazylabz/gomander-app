package usecases

import (
	"encoding/json"
	"fmt"
	"io/fs"

	"gomander/internal/localization/domain"
)

type GetTranslation interface {
	Execute(locale string) (*domain.Localization, error)
}

type DefaultGetTranslation struct {
	localeFs fs.FS
}

func NewGetTranslation(localeFs fs.FS) *DefaultGetTranslation {
	return &DefaultGetTranslation{
		localeFs: localeFs,
	}
}

func (uc *DefaultGetTranslation) Execute(locale string) (*domain.Localization, error) {
	localeJson, err := fs.ReadFile(uc.localeFs, fmt.Sprintf("locales/%s.json", locale))
	if err != nil {
		return nil, fmt.Errorf("read locale json: %w", err)
	}

	var lng domain.Localization
	if err := json.Unmarshal(localeJson, &lng); err != nil {
		return nil, fmt.Errorf("unmarshal locale json: %w", err)
	}

	return &lng, nil
}