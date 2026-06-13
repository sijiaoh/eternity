package game

import (
	"strings"
	"testing"
	"testing/fstest"

	"eternity/internal/i18n"
	"eternity/internal/scenario"
)

// resolveScene is the launcher's selection logic. These specs pin the three behaviors the
// -scene flag promises: a name selects its scene, no name falls back to the title screen,
// and an unknown name fails loudly instead of booting somewhere unexpected.
func TestResolveScene(t *testing.T) {
	// Deliberately unsorted to confirm selection doesn't rely on ordering.
	available := []string{sceneBattle, sceneTitle}

	t.Run("a known name selects that scene", func(t *testing.T) {
		for _, name := range available {
			got, err := resolveScene(name, DefaultScene, available)
			if err != nil {
				t.Fatalf("resolveScene(%q) returned error: %v", name, err)
			}
			if got != name {
				t.Fatalf("resolveScene(%q) = %q, want %q", name, got, name)
			}
		}
	})

	t.Run("no name falls back to the title screen", func(t *testing.T) {
		got, err := resolveScene("", DefaultScene, available)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != sceneTitle {
			t.Fatalf("default launch = %q, want title screen %q", got, sceneTitle)
		}
	})

	t.Run("an unknown name is rejected", func(t *testing.T) {
		_, err := resolveScene("nope", DefaultScene, available)
		if err == nil {
			t.Fatal("expected an error for an unknown scene name")
		}
		// The error must name the bad request and list the valid options so a developer can
		// fix the flag without reading the source.
		if !strings.Contains(err.Error(), "nope") {
			t.Errorf("error %q should mention the unknown name", err)
		}
		for _, name := range available {
			if !strings.Contains(err.Error(), name) {
				t.Errorf("error %q should list valid scene %q", err, name)
			}
		}
	})
}

// validateSituation guards the launcher promise that a -scenario file reproduces an exact start
// state: a battle situation must target the battle scene, or the file silently boots elsewhere.
func TestValidateSituation(t *testing.T) {
	battle := scenario.Battle{Dialogue: true}

	t.Run("battle situation with the battle scene is accepted", func(t *testing.T) {
		if err := validateSituation(sceneBattle, battle); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("a non-battle scene with no situation is accepted", func(t *testing.T) {
		if err := validateSituation(sceneTitle, scenario.Battle{}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("battle situation with a non-battle scene is rejected", func(t *testing.T) {
		err := validateSituation(sceneTitle, battle)
		if err == nil {
			t.Fatal("expected an error when a battle situation targets a non-battle scene")
		}
		if !strings.Contains(err.Error(), sceneTitle) || !strings.Contains(err.Error(), sceneBattle) {
			t.Errorf("error %q should name both the wrong scene and the expected one", err)
		}
	})
}

// applyLocale is scene-agnostic: it sets the language before any scene is built, so these specs
// pin that an empty locale leaves the default, a known one is applied, and an unknown one fails
// loudly with the valid options rather than silently rendering the default language.
func TestApplyLocale(t *testing.T) {
	newBundle := func(t *testing.T) *i18n.Bundle {
		t.Helper()
		b, err := i18n.Load(fstest.MapFS{
			"zh.json": {Data: []byte(`{"k":"中"}`)},
			"en.json": {Data: []byte(`{"k":"en"}`)},
		}, "zh")
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		return b
	}

	t.Run("empty locale leaves the bundle default", func(t *testing.T) {
		b := newBundle(t)
		if err := applyLocale(b, ""); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if b.Locale() != "zh" {
			t.Fatalf("locale = %q, want default zh", b.Locale())
		}
	})

	t.Run("a known locale is applied", func(t *testing.T) {
		b := newBundle(t)
		if err := applyLocale(b, "en"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if b.Locale() != "en" {
			t.Fatalf("locale = %q, want en", b.Locale())
		}
	})

	t.Run("an unknown locale errors with the valid options", func(t *testing.T) {
		b := newBundle(t)
		err := applyLocale(b, "fr")
		if err == nil {
			t.Fatal("expected an error for an unknown locale")
		}
		if !strings.Contains(err.Error(), "fr") {
			t.Errorf("error %q should name the bad locale", err)
		}
		for _, l := range []string{"zh", "en"} {
			if !strings.Contains(err.Error(), l) {
				t.Errorf("error %q should list available locale %q", err, l)
			}
		}
	})
}
