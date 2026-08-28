package usecases

import (
	"encoding/json"
	"fmt"
	"io/fs"

	"gomander/internal/localization/domain"
)

type GetTranslation struct {
	localeFs fs.FS
}

func NewGetTranslation(localeFs fs.FS) *GetTranslation {
	return &GetTranslation{
		localeFs: localeFs,
	}
}

func (uc *GetTranslation) Execute(locale string) (*domain.Localization, error) {
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
