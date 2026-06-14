package system

import (
	"math"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

// FacingSystem updates facing direction based on velocity.
type FacingSystem struct {
	facings    *ecs.Storage[component.Facing]
	velocities *ecs.Storage[component.Velocity]
}

func NewFacingSystem(
	facings *ecs.Storage[component.Facing],
	velocities *ecs.Storage[component.Velocity],
) *FacingSystem {
	return &FacingSystem{
		facings:    facings,
		velocities: velocities,
	}
}

func (s *FacingSystem) Update(w *ecs.World, dt float64) {
	s.facings.Each(func(e ecs.Entity, facing *component.Facing) {
		if !w.Alive(e) {
			return
		}
		vel := s.velocities.GetPtr(e)
		if vel == nil {
			return
		}

		facing.Walking = vel.X != 0 || vel.Y != 0
		if !facing.Walking {
			// No movement: keep the current facing direction.
			return
		}

		// Symmetric 45° sectors: each axis direction owns its ±45° arc, so the
		// dominant velocity component decides facing. Diagonals resolve to the
		// nearer axis instead of always snapping horizontal.
		//
		// On the diagonal (|vel.X| == |vel.Y|) tie-break to horizontal: a fixed,
		// predictable rule matching the project's prior horizontal-priority feel.
		if math.Abs(vel.X) >= math.Abs(vel.Y) {
			if vel.X < 0 {
				facing.Direction = component.FacingLeft
			} else {
				facing.Direction = component.FacingRight
			}
		} else {
			if vel.Y < 0 {
				facing.Direction = component.FacingUp
			} else {
				facing.Direction = component.FacingDown
			}
		}
	})
}
