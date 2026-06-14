package system

import (
	"math"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

// RangedAISystem steers each RangedAI ally straight away from its nearest enemy at the ally's
// speed, so a ranged fighter keeps backing off as foes approach. It sets only the flee velocity —
// the cap that keeps the ally near the player is the Leash's job (see component.RangedAI). Enemy
// selection reuses nearestEnemy, the same foe-picking the chasers use. Run before MovementSystem.
type RangedAISystem struct {
	rangedAIs  *ecs.Storage[component.RangedAI]
	factions   *ecs.Storage[component.Faction]
	positions  *ecs.Storage[component.Position]
	velocities *ecs.Storage[component.Velocity]
}

func NewRangedAISystem(
	rangedAIs *ecs.Storage[component.RangedAI],
	factions *ecs.Storage[component.Faction],
	positions *ecs.Storage[component.Position],
	velocities *ecs.Storage[component.Velocity],
) *RangedAISystem {
	return &RangedAISystem{
		rangedAIs:  rangedAIs,
		factions:   factions,
		positions:  positions,
		velocities: velocities,
	}
}

func (s *RangedAISystem) Update(w *ecs.World, dt float64) {
	s.rangedAIs.Each(func(e ecs.Entity, ai *component.RangedAI) {
		if !w.Alive(e) {
			return
		}

		pos := s.positions.GetPtr(e)
		vel := s.velocities.GetPtr(e)
		if pos == nil || vel == nil {
			return
		}

		enemyPos, ok := s.positions.Get(nearestEnemy(w, e, s.factions, s.positions))
		if !ok {
			// No enemy to flee: hold position.
			vel.X = 0
			vel.Y = 0
			return
		}

		// Flee straight away from the enemy: direction (self - enemy), normalized × speed.
		dx := pos.X - enemyPos.X
		dy := pos.Y - enemyPos.Y
		dist := math.Hypot(dx, dy)
		if dist < 0.01 {
			// Enemy is on top of us: no well-defined away direction, so don't jitter.
			vel.X = 0
			vel.Y = 0
			return
		}

		vel.X = (dx / dist) * ai.Speed
		vel.Y = (dy / dist) * ai.Speed
	})
}
