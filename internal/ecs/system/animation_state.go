package system

import (
	"eternity/internal/component"
	"eternity/internal/ecs"
)

// AnimationStateSystem drives each entity's animation state from its facing. It does no
// direction reasoning itself: DirectionalAnimation resolves the facing into a state name
// (handling the four-vs-horizontal-only distinction) and Animation switches to it.
type AnimationStateSystem struct {
	animations   *ecs.Storage[component.Animation]
	facings      *ecs.Storage[component.Facing]
	directionals *ecs.Storage[component.DirectionalAnimation]
}

func NewAnimationStateSystem(
	animations *ecs.Storage[component.Animation],
	facings *ecs.Storage[component.Facing],
	directionals *ecs.Storage[component.DirectionalAnimation],
) *AnimationStateSystem {
	return &AnimationStateSystem{
		animations:   animations,
		facings:      facings,
		directionals: directionals,
	}
}

func (s *AnimationStateSystem) Update(w *ecs.World, dt float64) {
	s.facings.Each(func(e ecs.Entity, facing *component.Facing) {
		if !w.Alive(e) {
			return
		}
		anim := s.animations.GetPtr(e)
		dir := s.directionals.GetPtr(e)
		if anim == nil || dir == nil {
			return
		}

		anim.SetState(dir.ResolveState(facing.Walking, facing.Direction))
	})
}
