// Package scenario describes a debug launch: which scene to boot into and the initial
// situation to reproduce inside it. It is the format layer for the -scenario file (see
// docs/ebiten-code-rules.md「调试启动」) and holds only plain data, so it stays free of ebiten
// and is testable headless.
//
// Scenarios are for manual debugging and testing only; they never affect the normal player
// launch, which passes neither -scene nor -scenario.
package scenario

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Scenario is one debug launch configuration, mirroring VS Code's launch.json idea: a single
// self-contained file reproduces an exact start state. Scene names the boot scene (resolved by
// the launcher; empty falls back to the default player path). Locale is scene-agnostic: it sets
// the language before any scene is built, so any scene (title included) can be reproduced in a
// given language. Battle carries the situation for the battle scene and is ignored by scenes that
// have no situation to vary.
type Scenario struct {
	Scene  string `json:"scene"`
	Locale string `json:"locale,omitempty"` // initial locale; default: "" leaves the bundle default
	Battle Battle `json:"battle"`
}

// Battle parameterizes the battle scene's initial situation. Every field is optional: a missing
// field keeps the scene's normal default (see the resolvers in internal/scene), so a partial
// file only overrides what it names. Pointers distinguish "unset" from a meaningful zero (e.g.
// TimeScale 0 = paused, Goblin false = suppress the goblin).
type Battle struct {
	PlayerX   *float64 `json:"playerX,omitempty"`   // player start X in world units; default: screen center
	PlayerY   *float64 `json:"playerY,omitempty"`   // player start Y in world units; default: screen center
	Goblin    *bool    `json:"goblin,omitempty"`    // spawn the goblin; default: true
	GoblinX   *float64 `json:"goblinX,omitempty"`   // goblin start X in world units; default: offset from player
	GoblinY   *float64 `json:"goblinY,omitempty"`   // goblin start Y in world units; default: offset from player
	Dialogue  bool     `json:"dialogue,omitempty"`  // start already in the sample dialogue; default: false
	TimeScale *float64 `json:"timeScale,omitempty"` // initial time scale (0 = paused, 0.5 = slow-mo); default: 1.0
}

// IsZero reports whether no situation field is set — i.e. the battle scene would boot exactly as
// in normal play. Used to reject a scenario that asks for a situation the chosen scene can't apply.
func (b Battle) IsZero() bool {
	return b == Battle{}
}

// Parse decodes a scenario file. Unknown fields are rejected so a typo'd key fails loudly
// rather than being silently ignored — a debug file that doesn't do what it says is worse than
// no file. The error names the offending input so it can be fixed without reading the source.
func Parse(data []byte) (Scenario, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	var s Scenario
	if err := dec.Decode(&s); err != nil {
		return Scenario{}, fmt.Errorf("parse scenario: %w", err)
	}
	return s, nil
}
