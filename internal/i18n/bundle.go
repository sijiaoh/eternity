// Package i18n provides key-based text translation backed by per-locale data files,
// decoupling user-visible strings from code. It is deliberately small: a flat
// key→text lookup with default-locale fallback, which is all the game's dialogue and
// UI text require. The package has no ebiten dependency, so it is unit-testable under
// -tags test.
package i18n

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

// Bundle holds every locale's messages plus the default and currently-selected locale.
// Translations are keyed by a dotted namespace (e.g. "window.title").
type Bundle struct {
	messages      map[string]map[string]string // locale → key → text
	defaultLocale string
	current       string
}

// Load reads every "<locale>.json" file at the root of fsys into a Bundle. Each file is a
// flat JSON object of key→text whose name (sans ".json") is the locale code. defaultLocale
// is both the initial selection and the fallback locale, so it must have a data file.
func Load(fsys fs.FS, defaultLocale string) (*Bundle, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, err
	}

	messages := make(map[string]map[string]string)
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".json") {
			continue
		}
		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, err
		}
		var msgs map[string]string
		if err := json.Unmarshal(data, &msgs); err != nil {
			return nil, fmt.Errorf("i18n: parse %s: %w", name, err)
		}
		messages[strings.TrimSuffix(name, ".json")] = msgs
	}

	if _, ok := messages[defaultLocale]; !ok {
		return nil, fmt.Errorf("i18n: default locale %q has no data file", defaultLocale)
	}

	return &Bundle{
		messages:      messages,
		defaultLocale: defaultLocale,
		current:       defaultLocale,
	}, nil
}

// Get returns the text for key in the current locale, falling back to the default locale
// when the current locale lacks the key, and finally to the key itself so a missing
// translation surfaces visibly instead of rendering as blank.
func (b *Bundle) Get(key string) string {
	if text, ok := b.messages[b.current][key]; ok {
		return text
	}
	if text, ok := b.messages[b.defaultLocale][key]; ok {
		return text
	}
	return key
}

// SetLocale switches the active locale. It reports false and keeps the current locale when
// the locale has no loaded data, letting callers reject unknown selections.
func (b *Bundle) SetLocale(locale string) bool {
	if _, ok := b.messages[locale]; !ok {
		return false
	}
	b.current = locale
	return true
}

// Locale returns the currently active locale.
func (b *Bundle) Locale() string { return b.current }

// Locales returns the available locale codes in sorted order.
func (b *Bundle) Locales() []string {
	out := make([]string, 0, len(b.messages))
	for locale := range b.messages {
		out = append(out, locale)
	}
	sort.Strings(out)
	return out
}
