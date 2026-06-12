package system

import (
	"eternity/internal/component"
	"eternity/internal/ecs"
)

// AnimationStateSystem updates animation state based on facing direction.
type AnimationStateSystem struct {
	animations *ecs.Storage[component.Animation]
	facings    *ecs.Storage[component.Facing]
}

func NewAnimationStateSystem(
	animations *ecs.Storage[component.Animation],
	facings *ecs.Storage[component.Facing],
) *AnimationStateSystem {
	return &AnimationStateSystem{
		animations: animations,
		facings:    facings,
	}
}

func (s *AnimationStateSystem) Update(w *ecs.World, dt float64) {
	s.facings.Each(func(e ecs.Entity, facing *component.Facing) {
		if !w.Alive(e) {
			return
		}
		anim := s.animations.GetPtr(e)
		if anim == nil {
			return
		}

		stateName := facing.AnimationStateName()
		anim.SetState(stateName)
	})
}
