package infrastructure

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"

	"gomander/internal/localization"
)

const (
	localesDir      = "locales"
	localeExtension = ".json"
)

// LocaleFiles is the Catalogue as Gomander ships it: one JSON file per locale,
// packaged into the binary. The two constants above are the whole convention -
// a locale exists because a file is named after it, so adding a language is
// adding a file and nothing else.
type LocaleFiles struct {
	files fs.FS
}

func NewLocaleFiles(files fs.FS) *LocaleFiles {
	return &LocaleFiles{files: files}
}

func (l *LocaleFiles) Locales() ([]string, error) {
	dirEntries, err := fs.ReadDir(l.files, localesDir)
	if err != nil {
		return nil, fmt.Errorf("read locales directory: %w", err)
	}

	languages := make([]string, 0, len(dirEntries))
	for _, d := range dirEntries {
		if !d.IsDir() && strings.HasSuffix(d.Name(), localeExtension) {
			languageCode := strings.TrimSuffix(d.Name(), localeExtension)
			languages = append(languages, languageCode)
		}
	}

	return languages, nil
}

func (l *LocaleFiles) Translations(locale string) (*localization.Localization, error) {
	localeJson, err := fs.ReadFile(l.files, localesDir+"/"+locale+localeExtension)
	if err != nil {
		return nil, fmt.Errorf("read locale json: %w", err)
	}

	var lng localization.Localization
	if err := json.Unmarshal(localeJson, &lng); err != nil {
		return nil, fmt.Errorf("unmarshal locale json: %w", err)
	}

	return &lng, nil
}
