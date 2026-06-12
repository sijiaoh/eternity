//go:build !test

package system

import (
	"testing"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

func TestFacingSystem_UpdateDirection(t *testing.T) {
	w := ecs.NewWorld(10)
	facings := ecs.NewStorage[component.Facing](10)
	velocities := ecs.NewStorage[component.Velocity](10)

	sys := NewFacingSystem(facings, velocities)

	e := w.Spawn()
	facings.Set(e, component.Facing{Direction: component.FacingDown, Walking: false})
	velocities.Set(e, component.Velocity{X: -5, Y: 0})

	sys.Update(w, 0.016)

	facing, _ := facings.Get(e)
	if facing.Direction != component.FacingLeft {
		t.Errorf("direction = %v, want FacingLeft", facing.Direction)
	}
	if !facing.Walking {
		t.Error("walking should be true when velocity is non-zero")
	}
}

func TestFacingSystem_HorizontalPriority(t *testing.T) {
	w := ecs.NewWorld(10)
	facings := ecs.NewStorage[component.Facing](10)
	velocities := ecs.NewStorage[component.Velocity](10)

	sys := NewFacingSystem(facings, velocities)

	e := w.Spawn()
	facings.Set(e, component.Facing{Direction: component.FacingDown, Walking: false})
	velocities.Set(e, component.Velocity{X: 5, Y: -5}) // Moving right and up

	sys.Update(w, 0.016)

	facing, _ := facings.Get(e)
	// Horizontal (right) should take priority over vertical (up)
	if facing.Direction != component.FacingRight {
		t.Errorf("direction = %v, want FacingRight (horizontal priority)", facing.Direction)
	}
}

func TestFacingSystem_KeepsDirectionWhenStopped(t *testing.T) {
	w := ecs.NewWorld(10)
	facings := ecs.NewStorage[component.Facing](10)
	velocities := ecs.NewStorage[component.Velocity](10)

	sys := NewFacingSystem(facings, velocities)

	e := w.Spawn()
	facings.Set(e, component.Facing{Direction: component.FacingUp, Walking: true})
	velocities.Set(e, component.Velocity{X: 0, Y: 0})

	sys.Update(w, 0.016)

	facing, _ := facings.Get(e)
	// Direction should be preserved when stopped
	if facing.Direction != component.FacingUp {
		t.Errorf("direction = %v, want FacingUp (preserved)", facing.Direction)
	}
	if facing.Walking {
		t.Error("walking should be false when velocity is zero")
	}
}
