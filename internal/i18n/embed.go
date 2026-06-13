package i18n

import (
	"embed"
	"io/fs"
)

//go:embed locales/*.json
var embeddedLocales embed.FS

// DefaultLocale is the source language: all UI and dialogue text is authored in Chinese,
// so every other locale falls back to it.
const DefaultLocale = "zh"

// New loads the bundle from the locale files embedded in the binary, selecting
// DefaultLocale. It is the production entry point; tests use Load with a custom fs.FS.
func New() (*Bundle, error) {
	sub, err := fs.Sub(embeddedLocales, "locales")
	if err != nil {
		return nil, err
	}
	return Load(sub, DefaultLocale)
}
