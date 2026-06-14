package scene

import (
	"math"
	"testing"

	"eternity/internal/component"
	"eternity/internal/ecs"
	"eternity/internal/ecs/system"
)

// linkRangedCombat must arm every ally with the ranged behavior (speed derived from the player) and
// a leash to the anchor at the guard radius — the combat-state counterpart to the trail formation.
func TestLinkRangedCombat_ArmsAllies(t *testing.T) {
	w := ecs.NewWorld(10)
	rangedAIs := ecs.NewStorage[component.RangedAI](10)
	leashes := ecs.NewStorage[component.Leash](10)

	anchor := w.Spawn()
	first := w.Spawn()
	second := w.Spawn()

	linkRangedCombat(rangedAIs, leashes, anchor, []ecs.Entity{first, second})

	for _, ally := range []ecs.Entity{first, second} {
		ai, ok := rangedAIs.Get(ally)
		if !ok || ai.Speed != allySpeed() {
			t.Errorf("ally %v RangedAI = %+v (ok %v), want speed %v", ally, ai, ok, allySpeed())
		}
		leash, ok := leashes.Get(ally)
		if !ok || leash.Anchor != anchor || leash.Range != guardRange {
			t.Errorf("ally %v Leash = %+v (ok %v), want anchor %v range %v", ally, leash, ok, anchor, guardRange)
		}
	}

	// The anchor is just the reference point — it gets neither component.
	if _, ok := rangedAIs.Get(anchor); ok {
		t.Error("anchor should not receive a RangedAI")
	}
}

// updateParty must route to the right behavior: trailing in the normal state, ranged fleeing in the
// combat state. The leader and enemy sit so the two behaviors drive distinct directions, letting the
// assertions tell which system ran.
func TestUpdateParty_SwitchesBehaviorByState(t *testing.T) {
	w := ecs.NewWorld(10)
	positions := ecs.NewStorage[component.Position](10)
	velocities := ecs.NewStorage[component.Velocity](10)
	factions := ecs.NewStorage[component.Faction](10)
	trails := ecs.NewStorage[component.Trail](10)
	rangedAIs := ecs.NewStorage[component.RangedAI](10)

	leader := w.Spawn()
	positions.Set(leader, component.Position{X: 10, Y: 0}) // trailing heads +X toward the leader
	factions.Set(leader, component.FactionPlayer)

	enemy := w.Spawn()
	positions.Set(enemy, component.Position{X: 0, Y: 5}) // fleeing heads -Y away from the enemy
	factions.Set(enemy, component.FactionEnemy)

	ally := w.Spawn()
	positions.Set(ally, component.Position{X: 0, Y: 0})
	velocities.Set(ally, component.Velocity{})
	factions.Set(ally, component.FactionPlayer)
	trails.Set(ally, component.Trail{Leader: leader, Speed: 6, Gap: 1})
	rangedAIs.Set(ally, component.RangedAI{Speed: 6})

	trail := system.NewTrailSystem(trails, positions, velocities)
	ranged := system.NewRangedAISystem(rangedAIs, factions, positions, velocities)

	// Normal state: trail runs, so the ally heads toward the leader (+X).
	updateParty(false, trail, ranged, w, 0)
	vel, _ := velocities.Get(ally)
	if vel.X <= 0 || math.Abs(vel.Y) > 0.001 {
		t.Errorf("normal-state velocity = (%v, %v), want +X toward the leader (trailing)", vel.X, vel.Y)
	}

	// Combat state: ranged runs, so the ally flees the enemy (-Y).
	updateParty(true, trail, ranged, w, 0)
	vel, _ = velocities.Get(ally)
	if math.Abs(vel.X) > 0.001 || vel.Y >= 0 {
		t.Errorf("combat-state velocity = (%v, %v), want -Y away from the enemy (fleeing)", vel.X, vel.Y)
	}
}
