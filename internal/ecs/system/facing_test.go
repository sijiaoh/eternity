//go:build test

package system

import (
	"testing"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

// newFacingFixture wires the storages and system used across the facing tests.
func newFacingFixture() (*ecs.World, *ecs.Storage[component.Facing], *ecs.Storage[component.Velocity], *FacingSystem) {
	w := ecs.NewWorld(10)
	facings := ecs.NewStorage[component.Facing](10)
	velocities := ecs.NewStorage[component.Velocity](10)
	return w, facings, velocities, NewFacingSystem(facings, velocities)
}

func TestFacingSystem_Direction(t *testing.T) {
	// Y axis points down: vel.Y > 0 is Down, vel.Y < 0 is Up.
	tests := []struct {
		name string
		velX float64
		velY float64
		want component.FacingDirection
	}{
		// Pure axes.
		{"right", 5, 0, component.FacingRight},
		{"left", -5, 0, component.FacingLeft},
		{"down", 0, 5, component.FacingDown},
		{"up", 0, -5, component.FacingUp},

		// Diagonal quadrants resolve to the nearer axis (non-edge angles).
		{"right-dominant up", 5, -3, component.FacingRight},
		{"up-dominant right", 3, -5, component.FacingUp},
		{"left-dominant down", -5, 3, component.FacingLeft},
		{"down-dominant left", -3, 5, component.FacingDown},

		// 45° edges (|X| == |Y|): tie-break to horizontal, deterministically.
		{"edge right-up", 5, -5, component.FacingRight},
		{"edge right-down", 5, 5, component.FacingRight},
		{"edge left-up", -5, -5, component.FacingLeft},
		{"edge left-down", -5, 5, component.FacingLeft},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, facings, velocities, sys := newFacingFixture()

			// Seed a direction different from the expected one, so a pass proves
			// the system actively set the direction rather than leaving it as-is.
			seed := component.FacingDown
			if tt.want == seed {
				seed = component.FacingUp
			}
			e := w.Spawn()
			facings.Set(e, component.Facing{Direction: seed, Walking: false})
			velocities.Set(e, component.Velocity{X: tt.velX, Y: tt.velY})

			sys.Update(w, 0.016)

			facing, _ := facings.Get(e)
			if facing.Direction != tt.want {
				t.Errorf("direction = %v, want %v", facing.Direction, tt.want)
			}
			if !facing.Walking {
				t.Error("walking should be true when velocity is non-zero")
			}
		})
	}
}

func TestFacingSystem_KeepsDirectionWhenStopped(t *testing.T) {
	w, facings, velocities, sys := newFacingFixture()

	e := w.Spawn()
	facings.Set(e, component.Facing{Direction: component.FacingUp, Walking: true})
	velocities.Set(e, component.Velocity{X: 0, Y: 0})

	sys.Update(w, 0.016)

	facing, _ := facings.Get(e)
	if facing.Direction != component.FacingUp {
		t.Errorf("direction = %v, want FacingUp (preserved)", facing.Direction)
	}
	if facing.Walking {
		t.Error("walking should be false when velocity is zero")
	}
}
