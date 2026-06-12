//go:build !test

package system

import (
	"testing"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

func TestAnimationStateSystem_Update(t *testing.T) {
	w := ecs.NewWorld(10)
	animations := ecs.NewStorage[component.Animation](10)
	facings := ecs.NewStorage[component.Facing](10)

	sys := NewAnimationStateSystem(animations, facings)

	states := []component.AnimationState{
		{Name: "idle_down", StartFrame: 0, FrameCount: 1, FPS: 8, Loop: true},
		{Name: "walk_down", StartFrame: 0, FrameCount: 6, FPS: 8, Loop: true},
		{Name: "idle_left", StartFrame: 6, FrameCount: 1, FPS: 8, Loop: true},
		{Name: "walk_left", StartFrame: 6, FrameCount: 6, FPS: 8, Loop: true},
	}

	e := w.Spawn()
	animations.Set(e, *component.NewAnimation(states))
	facings.Set(e, component.Facing{Direction: component.FacingLeft, Walking: true})

	sys.Update(w, 0.016)

	anim, _ := animations.Get(e)
	if anim.State() != "walk_left" {
		t.Errorf("state = %s, want walk_left", anim.State())
	}
}

func TestAnimationStateSystem_IdleState(t *testing.T) {
	w := ecs.NewWorld(10)
	animations := ecs.NewStorage[component.Animation](10)
	facings := ecs.NewStorage[component.Facing](10)

	sys := NewAnimationStateSystem(animations, facings)

	states := []component.AnimationState{
		{Name: "idle_down", StartFrame: 0, FrameCount: 1, FPS: 8, Loop: true},
		{Name: "walk_down", StartFrame: 0, FrameCount: 6, FPS: 8, Loop: true},
	}

	e := w.Spawn()
	anim := component.NewAnimation(states)
	anim.SetState("walk_down") // Start in walking state
	animations.Set(e, *anim)
	facings.Set(e, component.Facing{Direction: component.FacingDown, Walking: false})

	sys.Update(w, 0.016)

	updatedAnim, _ := animations.Get(e)
	if updatedAnim.State() != "idle_down" {
		t.Errorf("state = %s, want idle_down", updatedAnim.State())
	}
}
