package scenario

import (
	"strings"
	"testing"
)

// Parse is the scenario file's contract. These specs pin what a debug file promises: named
// fields land on the struct, omitted fields stay unset (so the scene applies its defaults), and
// malformed or mistyped input fails loudly instead of launching a silently-wrong situation.
func TestParse(t *testing.T) {
	t.Run("a full file populates every field", func(t *testing.T) {
		s, err := Parse([]byte(`{
			"scene": "battle",
			"locale": "en",
			"battle": {
				"mageX": 1.5, "mageY": 2.5,
				"goblin": false,
				"goblinX": 7, "goblinY": 8,
				"dialogue": true,
				"timeScale": 0.5
			}
		}`))
		if err != nil {
			t.Fatalf("Parse: %v", err)
		}
		if s.Scene != "battle" {
			t.Errorf("Scene = %q, want battle", s.Scene)
		}
		if s.Locale != "en" {
			t.Errorf("Locale = %q, want en", s.Locale)
		}
		b := s.Battle
		if got := derefF(t, b.MageX); got != 1.5 {
			t.Errorf("MageX = %v, want 1.5", got)
		}
		if got := derefF(t, b.MageY); got != 2.5 {
			t.Errorf("MageY = %v, want 2.5", got)
		}
		if b.Goblin == nil || *b.Goblin {
			t.Errorf("Goblin = %v, want explicit false", b.Goblin)
		}
		if got := derefF(t, b.GoblinX); got != 7 {
			t.Errorf("GoblinX = %v, want 7", got)
		}
		if got := derefF(t, b.GoblinY); got != 8 {
			t.Errorf("GoblinY = %v, want 8", got)
		}
		if !b.Dialogue {
			t.Error("Dialogue = false, want true")
		}
		if got := derefF(t, b.TimeScale); got != 0.5 {
			t.Errorf("TimeScale = %v, want 0.5", got)
		}
	})

	t.Run("omitted fields stay unset for the scene to default", func(t *testing.T) {
		s, err := Parse([]byte(`{"scene": "battle"}`))
		if err != nil {
			t.Fatalf("Parse: %v", err)
		}
		b := s.Battle
		if b.MageX != nil || b.MageY != nil || b.GoblinX != nil || b.GoblinY != nil {
			t.Error("position fields should be nil when omitted")
		}
		if b.Goblin != nil {
			t.Error("Goblin should be nil when omitted")
		}
		if b.TimeScale != nil {
			t.Error("TimeScale should be nil when omitted")
		}
		if b.Dialogue {
			t.Error("Dialogue should be false when omitted")
		}
		if s.Locale != "" {
			t.Error("Locale should be empty when omitted")
		}
	})

	t.Run("a top-level locale leaves the battle situation empty", func(t *testing.T) {
		// The point of a scene-agnostic locale: reproduce a non-battle scene (e.g. the title) in a
		// chosen language. Locale must stay out of Battle so the situation reads as empty and the
		// launcher's "battle situation needs the battle scene" gate doesn't reject a title+locale file.
		s, err := Parse([]byte(`{"scene": "title", "locale": "en"}`))
		if err != nil {
			t.Fatalf("Parse: %v", err)
		}
		if s.Locale != "en" {
			t.Errorf("Locale = %q, want en", s.Locale)
		}
		if !s.Battle.IsZero() {
			t.Error("a top-level locale should not populate the battle situation")
		}
	})

	t.Run("the old battle.locale location is rejected after moving locale to the top level", func(t *testing.T) {
		// locale used to live under battle; it now belongs at the top level (scene-agnostic).
		// DisallowUnknownFields turns a stale battle.locale into a loud error rather than a
		// silently-ignored field, so old debug files fail instead of booting the wrong language.
		_, err := Parse([]byte(`{"scene": "battle", "battle": {"locale": "en"}}`))
		if err == nil {
			t.Fatal("expected an error for a locale nested under battle")
		}
		if !strings.Contains(err.Error(), "locale") {
			t.Errorf("error %q should name the offending field", err)
		}
	})

	t.Run("malformed JSON is rejected", func(t *testing.T) {
		if _, err := Parse([]byte(`{not json`)); err == nil {
			t.Fatal("expected an error for malformed JSON")
		}
	})

	t.Run("an unknown field is rejected so typos fail loud", func(t *testing.T) {
		_, err := Parse([]byte(`{"scene": "battle", "battle": {"goblnn": false}}`))
		if err == nil {
			t.Fatal("expected an error for an unknown field")
		}
		if !strings.Contains(err.Error(), "goblnn") {
			t.Errorf("error %q should name the unknown field", err)
		}
	})
}

// IsZero gates the launcher's "battle situation needs the battle scene" check, so it must report
// true only when nothing is set (normal play) and false the moment any field is.
func TestBattleIsZero(t *testing.T) {
	if !(Battle{}).IsZero() {
		t.Error("an unset Battle should be zero")
	}
	x := 1.0
	for name, b := range map[string]Battle{
		"position":  {MageX: &x},
		"goblin":    {Goblin: new(bool)},
		"dialogue":  {Dialogue: true},
		"timeScale": {TimeScale: &x},
	} {
		if b.IsZero() {
			t.Errorf("Battle with %s set should not be zero", name)
		}
	}
}

func derefF(t *testing.T, p *float64) float64 {
	t.Helper()
	if p == nil {
		t.Fatal("expected a value, got nil")
	}
	return *p
}
