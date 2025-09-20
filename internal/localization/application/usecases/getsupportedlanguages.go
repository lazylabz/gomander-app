package usecases

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

type GetSupportedLanguages interface {
	Execute() ([]string, error)
}

type DefaultGetSupportedLanguages struct {
	localeFs embed.FS
}

func NewGetSupportedLanguages(localeFs embed.FS) *DefaultGetSupportedLanguages {
	return &DefaultGetSupportedLanguages{
		localeFs: localeFs,
	}
}

func (uc *DefaultGetSupportedLanguages) Execute() ([]string, error) {
	dirEntries, err := fs.ReadDir(uc.localeFs, "locales")
	if err != nil {
		return nil, fmt.Errorf("read locales directory: %w", err)
	}

	languages := make([]string, 0, len(dirEntries))
	for _, d := range dirEntries {
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".json") {
			languageCode := strings.TrimSuffix(d.Name(), ".json")
			languages = append(languages, languageCode)
		}
	}

	return languages, nil
}