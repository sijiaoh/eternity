package system

import (
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

		// Priority: horizontal over vertical (common in 2D games)
		if vel.X < 0 {
			facing.Direction = component.FacingLeft
		} else if vel.X > 0 {
			facing.Direction = component.FacingRight
		} else if vel.Y < 0 {
			facing.Direction = component.FacingUp
		} else if vel.Y > 0 {
			facing.Direction = component.FacingDown
		}
		// If no movement, keep the current facing direction
	})
}
