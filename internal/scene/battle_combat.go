package scene

import (
	"eternity/internal/component"
	"eternity/internal/ecs"
	"eternity/internal/ecs/system"
)

// guardRange is the radius, in units, of the zone the party guards around the player. It is the
// single source of truth for two coupled rules: combat begins when an enemy enters this zone, and
// during combat allies are leashed within it. Tying both to one value keeps "the guard zone" a
// single concept — the enemy that starts the fight is exactly the one that breached the ring the
// allies hold. 5 units is an arbitrary but readable size (story: "随便设个5unit").
const guardRange = 5.0

// linkRangedCombat arms each ally with the ranged-combat behavior: flee the nearest enemy
// (RangedAI) while leashed within guardRange of the anchor (the player mage). It is the combat-state
// counterpart to linkTrailFormation; the scene gives an ally both, and updateParty picks which one
// drives it each frame.
func linkRangedCombat(rangedAIs *ecs.Storage[component.RangedAI], leashes *ecs.Storage[component.Leash], anchor ecs.Entity, allies []ecs.Entity) {
	for _, a := range allies {
		rangedAIs.Set(a, component.RangedAI{Speed: allySpeed()})
		leashes.Set(a, component.Leash{Anchor: anchor, Range: guardRange})
	}
}

// updateParty runs the allies' active behavior for the frame: their per-AI combat behavior when the
// party is in combat, otherwise the normal trailing. Keeping the normal/combat switch — the heart
// of the state change — in this build-tag-free file lets it be tested headlessly.
func updateParty(inCombat bool, trail *system.TrailSystem, ranged *system.RangedAISystem, w *ecs.World, dt float64) {
	if inCombat {
		ranged.Update(w, dt)
	} else {
		trail.Update(w, dt)
	}
}
