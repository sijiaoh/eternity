package system

import (
	"math"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

// AIFollowSystem updates velocity for entities following a target.
type AIFollowSystem struct {
	aiFollows  *ecs.Storage[component.AIFollow]
	positions  *ecs.Storage[component.Position]
	velocities *ecs.Storage[component.Velocity]
}

func NewAIFollowSystem(
	aiFollows *ecs.Storage[component.AIFollow],
	positions *ecs.Storage[component.Position],
	velocities *ecs.Storage[component.Velocity],
) *AIFollowSystem {
	return &AIFollowSystem{
		aiFollows:  aiFollows,
		positions:  positions,
		velocities: velocities,
	}
}

func (s *AIFollowSystem) Update(w *ecs.World, dt float64) {
	s.aiFollows.Each(func(e ecs.Entity, ai *component.AIFollow) {
		if !w.Alive(e) {
			return
		}

		pos := s.positions.GetPtr(e)
		vel := s.velocities.GetPtr(e)
		if pos == nil || vel == nil {
			return
		}

		if !w.Alive(ai.Target) {
			vel.X = 0
			vel.Y = 0
			return
		}

		targetPos, ok := s.positions.Get(ai.Target)
		if !ok {
			vel.X = 0
			vel.Y = 0
			return
		}

		dx := targetPos.X - pos.X
		dy := targetPos.Y - pos.Y
		dist := math.Hypot(dx, dy)

		if dist < 0.01 {
			vel.X = 0
			vel.Y = 0
			return
		}

		vel.X = (dx / dist) * ai.Speed
		vel.Y = (dy / dist) * ai.Speed
	})
}
