//go:build test

package system

import (
	"testing"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

func TestMovementSystem_Update(t *testing.T) {
	w := ecs.NewWorld(10)
	positions := ecs.NewStorage[component.Position](10)
	velocities := ecs.NewStorage[component.Velocity](10)

	sys := NewMovementSystem(positions, velocities)

	e := w.Spawn()
	positions.Set(e, component.Position{X: 0, Y: 0})
	velocities.Set(e, component.Velocity{X: 5, Y: 10})

	sys.Update(w, 0.5)

	pos, _ := positions.Get(e)
	if pos.X != 2.5 || pos.Y != 5 {
		t.Errorf("position = (%v, %v), want (2.5, 5)", pos.X, pos.Y)
	}
}

func TestMovementSystem_SkipsDeadEntities(t *testing.T) {
	w := ecs.NewWorld(10)
	positions := ecs.NewStorage[component.Position](10)
	velocities := ecs.NewStorage[component.Velocity](10)

	sys := NewMovementSystem(positions, velocities)

	e := w.Spawn()
	positions.Set(e, component.Position{X: 0, Y: 0})
	velocities.Set(e, component.Velocity{X: 10, Y: 10})

	w.Despawn(e)
	sys.Update(w, 1.0)

	// Position should not change for despawned entity
	pos, _ := positions.Get(e)
	if pos.X != 0 || pos.Y != 0 {
		t.Errorf("despawned entity should not move, got (%v, %v)", pos.X, pos.Y)
	}
}

func TestMovementSystem_RequiresBothComponents(t *testing.T) {
	w := ecs.NewWorld(10)
	positions := ecs.NewStorage[component.Position](10)
	velocities := ecs.NewStorage[component.Velocity](10)

	sys := NewMovementSystem(positions, velocities)

	e := w.Spawn()
	velocities.Set(e, component.Velocity{X: 10, Y: 10})
	// No position component

	// Should not panic
	sys.Update(w, 1.0)
}
