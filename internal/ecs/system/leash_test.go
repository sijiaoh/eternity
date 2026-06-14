//go:build test

package system

import (
	"math"
	"testing"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

// newLeashFixture wires the storages and system with an anchor at the origin, returning the anchor
// so specs place the leashed unit relative to it.
func newLeashFixture() (*ecs.World, *ecs.Storage[component.Position], *ecs.Storage[component.Velocity], *ecs.Storage[component.Leash], *LeashSystem, ecs.Entity) {
	w := ecs.NewWorld(10)
	positions := ecs.NewStorage[component.Position](10)
	velocities := ecs.NewStorage[component.Velocity](10)
	leashes := ecs.NewStorage[component.Leash](10)

	anchor := w.Spawn()
	positions.Set(anchor, component.Position{X: 0, Y: 0})

	sys := NewLeashSystem(leashes, positions, velocities)
	return w, positions, velocities, leashes, sys, anchor
}

// Inside the radius the leash is slack: it touches neither position nor velocity.
func TestLeashSystem_LeavesUnitWithinRange(t *testing.T) {
	w, positions, velocities, leashes, sys, anchor := newLeashFixture()

	unit := w.Spawn()
	positions.Set(unit, component.Position{X: 3, Y: 0}) // within the radius of 5
	velocities.Set(unit, component.Velocity{X: 4, Y: 0})
	leashes.Set(unit, component.Leash{Anchor: anchor, Range: 5})

	sys.Update(w, 0)

	pos, _ := positions.Get(unit)
	vel, _ := velocities.Get(unit)
	if pos.X != 3 || pos.Y != 0 || vel.X != 4 || vel.Y != 0 {
		t.Errorf("within range got pos (%v,%v) vel (%v,%v), want pos (3,0) vel (4,0) untouched", pos.X, pos.Y, vel.X, vel.Y)
	}
}

// Beyond the radius the unit is snapped back onto the rim along the anchor→unit direction.
func TestLeashSystem_SnapsUnitOntoRim(t *testing.T) {
	w, positions, velocities, leashes, sys, anchor := newLeashFixture()

	unit := w.Spawn()
	positions.Set(unit, component.Position{X: 6, Y: 8}) // distance 10, well past the radius of 5
	velocities.Set(unit, component.Velocity{})
	leashes.Set(unit, component.Leash{Anchor: anchor, Range: 5})

	sys.Update(w, 0)

	// Direction (6,8) normalized (0.6,0.8) × range 5 = (3,4): same bearing, pulled to the rim.
	pos, _ := positions.Get(unit)
	if math.Abs(pos.X-3) > 0.001 || math.Abs(pos.Y-4) > 0.001 {
		t.Errorf("position = (%v, %v), want clamped to rim (3, 4)", pos.X, pos.Y)
	}
}

// At the rim, the outward part of the velocity is cancelled so the unit stops pushing past it:
// a purely outward push leaves it resting (zero), while a tangential (along-rim) push survives.
func TestLeashSystem_CancelsOutwardVelocityAtRim(t *testing.T) {
	w, positions, velocities, leashes, sys, anchor := newLeashFixture()

	unit := w.Spawn()
	positions.Set(unit, component.Position{X: 6, Y: 0}) // just past the radius, due +X
	// Velocity mixes outward (+X, gets cancelled) and tangential (+Y, kept).
	velocities.Set(unit, component.Velocity{X: 7, Y: 2})
	leashes.Set(unit, component.Leash{Anchor: anchor, Range: 5})

	sys.Update(w, 0)

	vel, _ := velocities.Get(unit)
	if math.Abs(vel.X) > 0.001 || math.Abs(vel.Y-2) > 0.001 {
		t.Errorf("velocity = (%v, %v), want (0, 2): outward cancelled, tangential kept", vel.X, vel.Y)
	}
}

// A dead anchor means no tether, so the unit is left where it is rather than snapped to a stale spot.
func TestLeashSystem_IgnoresDeadAnchor(t *testing.T) {
	w, positions, velocities, leashes, sys, anchor := newLeashFixture()

	unit := w.Spawn()
	positions.Set(unit, component.Position{X: 20, Y: 0})
	velocities.Set(unit, component.Velocity{})
	leashes.Set(unit, component.Leash{Anchor: anchor, Range: 5})

	w.Despawn(anchor)
	sys.Update(w, 0)

	pos, _ := positions.Get(unit)
	if pos.X != 20 || pos.Y != 0 {
		t.Errorf("position with dead anchor = (%v, %v), want unchanged (20, 0)", pos.X, pos.Y)
	}
}
