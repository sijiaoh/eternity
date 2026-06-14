package system

import (
	"math"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

// TrailSystem moves each Trail entity toward its leader, halting once within the leader's Gap so
// the follower settles behind it instead of overlapping. Chained leaders form a single-file line
// (see component.Trail). It reads the leader's current position, so run it before MovementSystem.
type TrailSystem struct {
	trails     *ecs.Storage[component.Trail]
	positions  *ecs.Storage[component.Position]
	velocities *ecs.Storage[component.Velocity]
}

func NewTrailSystem(
	trails *ecs.Storage[component.Trail],
	positions *ecs.Storage[component.Position],
	velocities *ecs.Storage[component.Velocity],
) *TrailSystem {
	return &TrailSystem{
		trails:     trails,
		positions:  positions,
		velocities: velocities,
	}
}

func (s *TrailSystem) Update(w *ecs.World, dt float64) {
	s.trails.Each(func(e ecs.Entity, t *component.Trail) {
		if !w.Alive(e) {
			return
		}

		pos := s.positions.GetPtr(e)
		vel := s.velocities.GetPtr(e)
		if pos == nil || vel == nil {
			return
		}

		if !w.Alive(t.Leader) {
			vel.X = 0
			vel.Y = 0
			return
		}

		leaderPos, ok := s.positions.Get(t.Leader)
		if !ok {
			vel.X = 0
			vel.Y = 0
			return
		}

		dx := leaderPos.X - pos.X
		dy := leaderPos.Y - pos.Y
		dist := math.Hypot(dx, dy)
		if dist <= t.Gap {
			vel.X = 0
			vel.Y = 0
			return
		}

		vel.X = (dx / dist) * t.Speed
		vel.Y = (dy / dist) * t.Speed
	})
}
