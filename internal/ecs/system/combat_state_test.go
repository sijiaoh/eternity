//go:build test

package system

import (
	"testing"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

// newCombatFixture wires the storages and a player-faction unit at the origin, returning the system
// (radius 5) and the player so specs only place enemies.
func newCombatFixture() (*ecs.World, *ecs.Storage[component.Position], *ecs.Storage[component.Faction], *CombatStateSystem, ecs.Entity) {
	w := ecs.NewWorld(10)
	positions := ecs.NewStorage[component.Position](10)
	factions := ecs.NewStorage[component.Faction](10)

	player := w.Spawn()
	positions.Set(player, component.Position{X: 0, Y: 0})
	factions.Set(player, component.FactionPlayer)

	sys := NewCombatStateSystem(factions, positions, player, 5)
	return w, positions, factions, sys, player
}

func placeEnemy(w *ecs.World, positions *ecs.Storage[component.Position], factions *ecs.Storage[component.Faction], x, y float64) ecs.Entity {
	e := w.Spawn()
	positions.Set(e, component.Position{X: x, Y: y})
	factions.Set(e, component.FactionEnemy)
	return e
}

// An enemy on the radius (distance == 5) counts as breaching the guard zone: the boundary is
// inclusive, so combat engages exactly when the enemy reaches it.
func TestCombatState_EnemyOnRadiusEngages(t *testing.T) {
	w, positions, factions, sys, _ := newCombatFixture()
	placeEnemy(w, positions, factions, 5, 0) // exactly on the radius

	if !sys.InCombat(w) {
		t.Error("InCombat = false with an enemy on the radius, want true")
	}
}

// Just outside the radius the party is still at peace — the boundary is the dividing line.
func TestCombatState_EnemyOutsideRadiusIsPeaceful(t *testing.T) {
	w, positions, factions, sys, _ := newCombatFixture()
	placeEnemy(w, positions, factions, 5.001, 0)

	if sys.InCombat(w) {
		t.Error("InCombat = true with the enemy just outside the radius, want false")
	}
}

// A same-faction unit close to the player must not trigger combat — only opposing factions count.
// This is also the "no qualifying enemy in range → peaceful" baseline: the scan finds a unit but
// none that counts, so the party stays at peace.
func TestCombatState_IgnoresSameFaction(t *testing.T) {
	w, positions, factions, sys, _ := newCombatFixture()
	friend := w.Spawn()
	positions.Set(friend, component.Position{X: 1, Y: 0})
	factions.Set(friend, component.FactionPlayer)

	if sys.InCombat(w) {
		t.Error("InCombat = true with only a same-faction unit nearby, want false")
	}
}

// With the player gone there is no party to command, so combat is never reported even with a
// nearby enemy.
func TestCombatState_DeadPlayerIsPeaceful(t *testing.T) {
	w, positions, factions, sys, player := newCombatFixture()
	placeEnemy(w, positions, factions, 1, 0)

	w.Despawn(player)
	if sys.InCombat(w) {
		t.Error("InCombat = true with a dead player, want false")
	}
}
