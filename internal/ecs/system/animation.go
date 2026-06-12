package system

import (
	"eternity/internal/component"
	"eternity/internal/ecs"
)

// AnimationSystem updates animation state for all entities with Animation component.
type AnimationSystem struct {
	animations *ecs.Storage[component.Animation]
}

func NewAnimationSystem(animations *ecs.Storage[component.Animation]) *AnimationSystem {
	return &AnimationSystem{animations: animations}
}

func (s *AnimationSystem) Update(w *ecs.World, dt float64) {
	s.animations.Each(func(e ecs.Entity, anim *component.Animation) {
		if !w.Alive(e) {
			return
		}
		anim.Update(dt)
	})
}
