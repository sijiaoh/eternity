package system

import (
	"math"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

// AITargetingSystem chooses each AIFollow entity's target: the nearest living unit of an
// opposing faction. It rewrites AIFollow.Target every frame, so the follower re-picks as
// enemies move, die, or appear. AIFollowSystem then moves toward whatever target is set and
// stops when there is none (an invalid target). Run this before AIFollowSystem.
type AITargetingSystem struct {
	aiFollows *ecs.Storage[component.AIFollow]
	factions  *ecs.Storage[component.Faction]
	positions *ecs.Storage[component.Position]
}

func NewAITargetingSystem(
	aiFollows *ecs.Storage[component.AIFollow],
	factions *ecs.Storage[component.Faction],
	positions *ecs.Storage[component.Position],
) *AITargetingSystem {
	return &AITargetingSystem{
		aiFollows: aiFollows,
		factions:  factions,
		positions: positions,
	}
}

func (s *AITargetingSystem) Update(w *ecs.World, dt float64) {
	s.aiFollows.Each(func(e ecs.Entity, ai *component.AIFollow) {
		if !w.Alive(e) {
			return
		}
		ai.Target = s.nearestEnemy(w, e)
	})
}

// nearestEnemy returns the closest living, positioned unit whose faction opposes e's, or the
// invalid zero entity when e has no faction or position, or no enemy exists.
func (s *AITargetingSystem) nearestEnemy(w *ecs.World, e ecs.Entity) ecs.Entity {
	myFaction, ok := s.factions.Get(e)
	if !ok {
		return ecs.Entity{}
	}
	pos, ok := s.positions.Get(e)
	if !ok {
		return ecs.Entity{}
	}

	var nearest ecs.Entity
	bestDist := math.MaxFloat64
	s.factions.Each(func(other ecs.Entity, faction *component.Faction) {
		if *faction == myFaction || !w.Alive(other) {
			return
		}
		otherPos, ok := s.positions.Get(other)
		if !ok {
			return
		}
		dist := math.Hypot(otherPos.X-pos.X, otherPos.Y-pos.Y)
		if dist < bestDist {
			bestDist = dist
			nearest = other
		}
	})
	return nearest
}
