package system

import (
	"eternity/internal/component"
	"eternity/internal/ecs"
)

// AnimationSystem advances animation playback for all entities with Animation component.
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
		updateAnimation(anim, dt)
	})
}

// updateAnimation advances the animation by deltaTime seconds.
func updateAnimation(anim *component.Animation, dt float64) {
	state, ok := anim.States[anim.CurrentState]
	if !ok || state.FrameCount <= 1 || state.FPS <= 0 || anim.Finished {
		return
	}
	if dt <= 0 {
		return
	}

	frameDuration := 1.0 / state.FPS
	anim.Elapsed += dt

	for anim.Elapsed >= frameDuration {
		anim.Elapsed -= frameDuration
		anim.FrameIndex++

		if anim.FrameIndex >= state.FrameCount {
			if state.Loop {
				anim.FrameIndex = 0
			} else {
				anim.FrameIndex = state.FrameCount - 1
				anim.Finished = true
				return
			}
		}
	}
}
