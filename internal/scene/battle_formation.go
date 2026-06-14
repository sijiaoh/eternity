package scene

import (
	"eternity/internal/component"
	"eternity/internal/ecs"
)

// Party movement tuning. These live in this build-tag-free file so the derivations stay
// headlessly testable, and so the player's speed has a single source feeding both the player's
// input control and the allies' trailing speed (DRY) rather than being written twice.
const (
	// playerSpeed is how fast the player-controlled mage moves, in units per second.
	playerSpeed = 5.0
	// allySpeedMultiplier makes allies a bit faster than the player so they can close the line
	// instead of falling behind.
	allySpeedMultiplier = 1.2
	// trailGap is how far behind its leader an ally settles, in units — about one body length, so
	// the line reads as single-file without the sprites overlapping.
	trailGap = 1.0
	// trailSlowRadius is the arrival-steering deceleration band past trailGap, in units. The ally
	// runs at full speed until it closes to within this band, then eases to a stop at the gap, so it
	// converges smoothly on the player instead of sprinting up and slamming to a halt.
	trailSlowRadius = 2.0
)

// allySpeed derives an ally's movement speed from the player's, expressing "allies move at 1.2× the
// player" as a relationship rather than a second hard-coded number. It drives both trailing and
// ranged-combat motion, so an ally moves at one speed regardless of state.
func allySpeed() float64 {
	return playerSpeed * allySpeedMultiplier
}

// linkTrailFormation arranges followers into a single-file line behind leader, Dragon-Quest style:
// the first follower trails the leader, each subsequent follower trails the one ahead of it. It
// scales to any number of allies even though the scene currently spawns one.
func linkTrailFormation(trails *ecs.Storage[component.Trail], leader ecs.Entity, followers []ecs.Entity) {
	ahead := leader
	for _, f := range followers {
		trails.Set(f, component.Trail{Leader: ahead, Speed: allySpeed(), Gap: trailGap, SlowRadius: trailSlowRadius})
		ahead = f
	}
}
