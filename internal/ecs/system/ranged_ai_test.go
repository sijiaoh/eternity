//go:build test

package system

import (
	"math"
	"testing"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

// newRangedFixture wires the storages and system used across the ranged-AI tests, with one
// player-faction ally at the origin so each spec only has to place enemies and run the system.
func newRangedFixture() (*ecs.World, *ecs.Storage[component.Position], *ecs.Storage[component.Velocity], *ecs.Storage[component.Faction], *RangedAISystem, ecs.Entity) {
	w := ecs.NewWorld(10)
	positions := ecs.NewStorage[component.Position](10)
	velocities := ecs.NewStorage[component.Velocity](10)
	factions := ecs.NewStorage[component.Faction](10)
	rangedAIs := ecs.NewStorage[component.RangedAI](10)

	ally := w.Spawn()
	positions.Set(ally, component.Position{X: 0, Y: 0})
	velocities.Set(ally, component.Velocity{})
	factions.Set(ally, component.FactionPlayer)
	rangedAIs.Set(ally, component.RangedAI{Speed: 10})

	sys := NewRangedAISystem(rangedAIs, factions, positions, velocities)
	return w, positions, velocities, factions, sys, ally
}

func addRangedEnemy(w *ecs.World, positions *ecs.Storage[component.Position], factions *ecs.Storage[component.Faction], x, y float64) ecs.Entity {
	e := w.Spawn()
	positions.Set(e, component.Position{X: x, Y: y})
	factions.Set(e, component.FactionEnemy)
	return e
}

// The ally backs straight away from its enemy at exactly Speed: direction (self - enemy) normalized
// (so distance doesn't inflate the velocity) and scaled by Speed.
func TestRangedAISystem_FleesDirectlyAwayAtSpeed(t *testing.T) {
	w, positions, velocities, factions, sys, ally := newRangedFixture()
	addRangedEnemy(w, positions, factions, -3, -4) // enemy down-left; away is up-right

	sys.Update(w, 0)

	// (self - enemy) = (3, 4), normalized (0.6, 0.8), scaled by speed 10 = (6, 8).
	vel, _ := velocities.Get(ally)
	if math.Abs(vel.X-6) > 0.001 || math.Abs(vel.Y-8) > 0.001 {
		t.Errorf("velocity = (%v, %v), want (6, 8) away from enemy", vel.X, vel.Y)
	}
}

// With several enemies the ally flees the nearest one, so closing foes drive the retreat.
func TestRangedAISystem_FleesNearestEnemy(t *testing.T) {
	w, positions, velocities, factions, sys, ally := newRangedFixture()
	addRangedEnemy(w, positions, factions, -10, 0) // far, to the left
	addRangedEnemy(w, positions, factions, 0, 2)   // near, above

	sys.Update(w, 0)

	// Nearest is above, so the ally flees straight down: (0, -speed).
	vel, _ := velocities.Get(ally)
	if math.Abs(vel.X) > 0.001 || math.Abs(vel.Y+10) > 0.001 {
		t.Errorf("velocity = (%v, %v), want (0, -10) away from the nearest enemy", vel.X, vel.Y)
	}
}

// No enemy on the field means nothing to flee, so the ally holds still.
func TestRangedAISystem_HoldsWhenNoEnemy(t *testing.T) {
	w, _, velocities, _, sys, ally := newRangedFixture()
	velocities.Set(ally, component.Velocity{X: 5, Y: 5})

	sys.Update(w, 0)

	vel, _ := velocities.Get(ally)
	if vel.X != 0 || vel.Y != 0 {
		t.Errorf("velocity with no enemy = (%v, %v), want (0, 0)", vel.X, vel.Y)
	}
}
