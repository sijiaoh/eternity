package scene

import (
	"eternity/internal/component"
	"eternity/internal/config"
	"eternity/internal/ecs"
	"eternity/internal/i18n"
	"eternity/internal/scenario"
)

// goblinSpawnOffset is how far the goblin spawns from the mage by default, in world units.
const goblinSpawnOffset = 3.0

// allySpawnOffset is how far the friendly mage stands to the player mage's left, in world units.
// It's fixed rather than scenario-driven — the ally has no debug knobs of its own.
const allySpawnOffset = 1.5

// The situation resolvers below turn an optional scenario.Battle into the concrete values
// NewBattleScene needs, applying the scene's normal defaults wherever a field is unset. They
// live in this build-tag-free file (rather than inline in the ebiten-gated battle.go) so the
// situation logic — the part with debug/test value — is unit-testable in a headless build.

// resolveMageStart returns the mage's start position in world units, defaulting to the
// screen center (the normal mage spawn) when the scenario leaves it unset.
func resolveMageStart(b scenario.Battle) (x, y float64) {
	cx := config.PixelsToUnits(config.ScreenWidth / 2)
	cy := config.PixelsToUnits(config.ScreenHeight / 2)
	return coalesce(b.MageX, cx), coalesce(b.MageY, cy)
}

// resolveGoblinStart returns the goblin's start position, defaulting to a fixed offset from the
// mage so the two spawn near each other as in normal play.
func resolveGoblinStart(b scenario.Battle, mageX, mageY float64) (x, y float64) {
	return coalesce(b.GoblinX, mageX+goblinSpawnOffset), coalesce(b.GoblinY, mageY+goblinSpawnOffset)
}

// resolveAllyStart returns the friendly mage's start position, a fixed offset to the player
// mage's left. With nothing driving its movement, this is also where the ally stays put.
func resolveAllyStart(mageX, mageY float64) (x, y float64) {
	return mageX - allySpawnOffset, mageY
}

// assignBattleFactions tags each unit with its side. Faction membership is a scene-layer
// decision (like InputControlled/CameraTarget). Player- and enemy-side entities are passed
// separately so the conditionally-spawned goblin is simply omitted from enemySide when absent.
func assignBattleFactions(factions *ecs.Storage[component.Faction], playerSide, enemySide []ecs.Entity) {
	for _, e := range playerSide {
		factions.Set(e, component.FactionPlayer)
	}
	for _, e := range enemySide {
		factions.Set(e, component.FactionEnemy)
	}
}

// spawnGoblin reports whether the goblin should be created. The default is true, matching normal
// play; a scenario can set goblin:false to debug the scene without it.
func spawnGoblin(b scenario.Battle) bool {
	return coalesce(b.Goblin, true)
}

// resolveTimeScale returns the initial clock scale, defaulting to 1.0 (normal speed). 0 pauses
// the world and values in between give slow-motion.
func resolveTimeScale(b scenario.Battle) float64 {
	return coalesce(b.TimeScale, 1.0)
}

// coalesce returns *p when set, else def — the one place the "unset means default" rule for
// optional scenario fields lives.
func coalesce[T any](p *T, def T) T {
	if p != nil {
		return *p
	}
	return def
}

// initialDialogue returns the lines the scene should open in, or nil when the scenario doesn't
// ask to start in dialogue. Passing nil to Dialogue.Start leaves it inactive, so the caller
// needs no conditional.
func initialDialogue(b scenario.Battle, bundle *i18n.Bundle) []component.DialogueLine {
	if !b.Dialogue {
		return nil
	}
	return sampleDialogue(bundle)
}

// sampleDialogue is the demo script for manual verification: a portrait line, a portrait-less
// narration line, and a long line that exercises CJK wrapping. Text and speaker name come from
// i18n, so the same script renders in whichever locale the bundle currently holds.
func sampleDialogue(b *i18n.Bundle) []component.DialogueLine {
	return []component.DialogueLine{
		{Speaker: b.Get("speaker.mage"), Portrait: "mage", Text: b.Get("dialogue.intro.warning")},
		{Speaker: "", Portrait: "", Text: b.Get("dialogue.intro.silence")},
		{Speaker: b.Get("speaker.mage"), Portrait: "mage", Text: b.Get("dialogue.intro.ready")},
	}
}
