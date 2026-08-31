// Package localization is the UI copy the frontend renders. A Localization is
// one locale's copy: a serialization format rather than an entity, since its
// tags are the i18n keys the frontend indexes and not a naming decision the
// backend gets to make. Which locales the Catalogue holds, how they are named
// and where they are packaged is the adapter's business alone.
package localization

// Catalogue is the translated copy Gomander ships, one Localization per locale.
type Catalogue interface {
	Locales() ([]string, error)
	Translations(locale string) (*Localization, error)
}
