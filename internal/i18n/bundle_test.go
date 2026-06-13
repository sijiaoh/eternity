package i18n

import (
	"fmt"
	"reflect"
	"testing"
	"testing/fstest"
)

// twoLocales builds an in-memory bundle where "en" is intentionally missing a key, so the
// fallback path is exercised without touching the real (deliberately complete) data files.
func twoLocales(t *testing.T) *Bundle {
	t.Helper()
	fsys := fstest.MapFS{
		"zh.json": {Data: []byte(`{"greeting":"你好","only_zh":"仅中文"}`)},
		"en.json": {Data: []byte(`{"greeting":"Hello"}`)},
	}
	b, err := Load(fsys, "zh")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return b
}

func TestGetUsesCurrentLocale(t *testing.T) {
	b := twoLocales(t)
	if !b.SetLocale("en") {
		t.Fatal("expected true when switching to a known locale")
	}
	if got := b.Get("greeting"); got != "Hello" {
		t.Fatalf("expected current-locale value, got %q", got)
	}
}

func TestGetFallsBackToDefaultLocale(t *testing.T) {
	b := twoLocales(t)
	b.SetLocale("en") // "only_zh" exists only in the default locale
	if got := b.Get("only_zh"); got != "仅中文" {
		t.Fatalf("expected fallback to default locale, got %q", got)
	}
}

func TestGetMissingKeyReturnsKey(t *testing.T) {
	b := twoLocales(t)
	if got := b.Get("nope"); got != "nope" {
		t.Fatalf("expected missing key echoed back, got %q", got)
	}
}

func TestSetLocaleUnknownKeepsCurrent(t *testing.T) {
	b := twoLocales(t)
	if b.SetLocale("fr") {
		t.Fatal("expected false for unknown locale")
	}
	if b.Locale() != "zh" {
		t.Fatalf("expected current locale unchanged, got %q", b.Locale())
	}
}

func TestLoadDefaultLocaleMustExist(t *testing.T) {
	fsys := fstest.MapFS{"en.json": {Data: []byte(`{}`)}}
	if _, err := Load(fsys, "zh"); err == nil {
		t.Fatal("expected error when default locale has no data file")
	}
}

func TestLocalesSorted(t *testing.T) {
	b := twoLocales(t)
	if got := b.Locales(); !reflect.DeepEqual(got, []string{"en", "zh"}) {
		t.Fatalf("expected sorted locales, got %v", got)
	}
}

// TestEmbeddedBundleLoads proves the embedded data files parse and the default locale is
// selected, guarding against a malformed JSON file slipping in.
func TestEmbeddedBundleLoads(t *testing.T) {
	b, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if b.Locale() != DefaultLocale {
		t.Fatalf("expected default locale %q, got %q", DefaultLocale, b.Locale())
	}
	if got := b.Get("window.title"); got == "window.title" {
		t.Fatal("expected window.title to be translated, got the key back")
	}
}

// TestEmbeddedLocalesComplete locks down the requirement that every shipped locale carries a
// full translation. The runtime fallback hides incompleteness: a non-default locale missing a
// key silently renders the Chinese default, so an English game could show stray Chinese with no
// error. It also flags keys absent from the default locale, catching typos and stale entries.
func TestEmbeddedLocalesComplete(t *testing.T) {
	b, err := New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	want := b.messages[DefaultLocale]
	for locale, msgs := range b.messages {
		if locale == DefaultLocale {
			continue
		}
		for key := range want {
			if _, ok := msgs[key]; !ok {
				t.Errorf("locale %q missing key %q present in default %q", locale, key, DefaultLocale)
			}
		}
		for key := range msgs {
			if _, ok := want[key]; !ok {
				t.Errorf("locale %q has key %q absent from default %q (stale or typo?)", locale, key, DefaultLocale)
			}
		}
	}
}

// ExampleBundle is the minimal runnable demo: load the embedded bundle, read a string, then
// switch locales and read the same key. Running it under `go test` self-proves the structure.
func ExampleBundle() {
	b, _ := New()

	fmt.Println(b.Get("speaker.mage")) // DefaultLocale ("zh")
	b.SetLocale("en")
	fmt.Println(b.Get("speaker.mage"))
	// Output:
	// 法师
	// Mage
}
