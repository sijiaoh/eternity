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

		stateName := facingAnimationStateName(facing)
		setAnimationState(anim, stateName)
	})
}

// facingAnimationStateName returns the animation state name based on facing and walking state.
func facingAnimationStateName(f *component.Facing) string {
	prefix := "idle_"
	if f.Walking {
		prefix = "walk_"
	}

	switch f.Direction {
	case component.FacingDown:
		return prefix + "down"
	case component.FacingLeft:
		return prefix + "left"
	case component.FacingRight:
		return prefix + "right"
	case component.FacingUp:
		return prefix + "up"
	default:
		return prefix + "down"
	}
}

// setAnimationState switches animation to a different state.
// If the state is the same as current, nothing happens.
// If the state doesn't exist, nothing happens.
func setAnimationState(anim *component.Animation, name string) {
	if anim.CurrentState == name {
		return
	}
	if _, ok := anim.States[name]; !ok {
		return
	}
	anim.CurrentState = name
	anim.FrameIndex = 0
	anim.Elapsed = 0
	anim.Finished = false
}
