package game

import (
	"fmt"
	"sort"
	"strings"

	"eternity/internal/i18n"
	"eternity/internal/scenario"
)

// Scene names are the launcher's vocabulary: each one keys a builder in New and is a valid
// value for the -scene flag. The string literals live here once so the rest of the code
// refers to scenes by constant.
const (
	sceneTitle  = "title"
	sceneBattle = "battle"

	// DefaultScene is the launch target when no scene is requested. It keeps the normal
	// player path unchanged: boot to the title screen.
	DefaultScene = sceneTitle
)

// resolveScene picks which registered scene to launch. An empty request falls back to
// defaultName so the default player path is unchanged; an unknown request is rejected with
// the list of valid names rather than silently booting somewhere unexpected.
//
// available is the set of registered scene names (the registry keys), passed in as a
// parameter so this selection logic stays free of the ebiten-dependent builders and is
// unit-testable in a headless build.
func resolveScene(requested, defaultName string, available []string) (string, error) {
	if requested == "" {
		requested = defaultName
	}
	for _, name := range available {
		if name == requested {
			return requested, nil
		}
	}

	valid := append([]string(nil), available...)
	sort.Strings(valid)
	return "", fmt.Errorf("unknown scene %q; available scenes: %s", requested, strings.Join(valid, ", "))
}

// validateSituation rejects a scenario whose situation can't be applied to the chosen scene.
// Only the battle scene has a situation to vary, so battle parameters with any other start scene
// is a mistake — most likely a -scenario file that set battle fields but forgot scene:"battle"
// (an omitted scene defaults to the title screen). Fail loud rather than silently boot the title
// and drop the situation the file promised to reproduce.
func validateSituation(scene string, battle scenario.Battle) error {
	if scene != sceneBattle && !battle.IsZero() {
		return fmt.Errorf("scenario sets a battle situation but start scene is %q; set scene to %q", scene, sceneBattle)
	}
	return nil
}

// applyLocale switches the bundle to the scenario's locale before any scene is built, so every
// scene (title included) resolves its text in the requested language. It is scene-agnostic and
// pure (no ebiten), so it stays unit-testable in a headless build. An empty locale leaves the
// bundle default untouched; an unknown one fails loudly with the valid options rather than
// silently rendering the default language.
func applyLocale(bundle *i18n.Bundle, locale string) error {
	if locale == "" {
		return nil
	}
	if !bundle.SetLocale(locale) {
		return fmt.Errorf("unknown locale %q; available: %v", locale, bundle.Locales())
	}
	return nil
}
