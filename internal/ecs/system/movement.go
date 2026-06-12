package system

import (
	"eternity/internal/component"
	"eternity/internal/ecs"
)

// MovementSystem applies velocity to position for all entities with both components.
type MovementSystem struct {
	positions  *ecs.Storage[component.Position]
	velocities *ecs.Storage[component.Velocity]
}

func NewMovementSystem(
	positions *ecs.Storage[component.Position],
	velocities *ecs.Storage[component.Velocity],
) *MovementSystem {
	return &MovementSystem{
		positions:  positions,
		velocities: velocities,
	}
}

func (s *MovementSystem) Update(w *ecs.World, dt float64) {
	s.velocities.Each(func(e ecs.Entity, vel *component.Velocity) {
		if !w.Alive(e) {
			return
		}
		pos := s.positions.GetPtr(e)
		if pos == nil {
			return
		}
		pos.X += vel.X * dt
		pos.Y += vel.Y * dt
	})
}
