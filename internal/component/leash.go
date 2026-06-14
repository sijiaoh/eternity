package component

import "eternity/internal/ecs"

// Leash tethers a unit to an Anchor, forbidding it from straying past Range units away. It is a
// position constraint enforced after movement by system.LeashSystem. Allies carry one so a combat
// AI (e.g. RangedAI fleeing an enemy) can never abandon the player it is anchored to.
type Leash struct {
	Anchor ecs.Entity // unit to stay near (the player mage)
	Range  float64    // farthest the unit may sit from the anchor, in units
}
