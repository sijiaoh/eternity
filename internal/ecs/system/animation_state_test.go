//go:build test

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
	if anim.CurrentState != "walk_left" {
		t.Errorf("state = %s, want walk_left", anim.CurrentState)
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
	// Start in walking state using setAnimationState (same logic the system uses)
	setAnimationState(anim, "walk_down")
	animations.Set(e, *anim)
	facings.Set(e, component.Facing{Direction: component.FacingDown, Walking: false})

	sys.Update(w, 0.016)

	updatedAnim, _ := animations.Get(e)
	if updatedAnim.CurrentState != "idle_down" {
		t.Errorf("state = %s, want idle_down", updatedAnim.CurrentState)
	}
}

func TestAnimationStateSystem_SetStateResetsAnimation(t *testing.T) {
	w := ecs.NewWorld(10)
	animations := ecs.NewStorage[component.Animation](10)
	facings := ecs.NewStorage[component.Facing](10)

	states := []component.AnimationState{
		{Name: "idle_down", StartFrame: 0, FrameCount: 6, FPS: 10, Loop: true},
		{Name: "walk_down", StartFrame: 6, FrameCount: 6, FPS: 12, Loop: true},
	}

	e := w.Spawn()
	anim := component.NewAnimation(states)
	// Advance the animation
	anim.Elapsed = 0.15
	anim.FrameIndex = 1
	animations.Set(e, *anim)
	facings.Set(e, component.Facing{Direction: component.FacingDown, Walking: true})

	sys := NewAnimationStateSystem(animations, facings)
	sys.Update(w, 0.016)

	animPtr := animations.GetPtr(e)
	if animPtr.CurrentState != "walk_down" {
		t.Errorf("state = %s, want walk_down", animPtr.CurrentState)
	}
	// Switching state should reset frame
	if animPtr.FrameIndex != 0 {
		t.Errorf("FrameIndex = %d, want 0 after state change", animPtr.FrameIndex)
	}
	if animPtr.Elapsed != 0 {
		t.Errorf("Elapsed = %f, want 0 after state change", animPtr.Elapsed)
	}
}

func TestAnimationStateSystem_SameStateNoReset(t *testing.T) {
	w := ecs.NewWorld(10)
	animations := ecs.NewStorage[component.Animation](10)
	facings := ecs.NewStorage[component.Facing](10)

	states := []component.AnimationState{
		{Name: "walk_down", StartFrame: 0, FrameCount: 6, FPS: 10, Loop: true},
	}

	e := w.Spawn()
	anim := component.NewAnimation(states)
	animations.Set(e, *anim)
	facings.Set(e, component.Facing{Direction: component.FacingDown, Walking: true})

	sys := NewAnimationStateSystem(animations, facings)

	// First update - set state to walk_down (same as initial)
	sys.Update(w, 0.016)

	// Manually advance the animation state
	animPtr := animations.GetPtr(e)
	animPtr.FrameIndex = 3
	animPtr.Elapsed = 0.1

	// Second update - same facing, should not reset
	sys.Update(w, 0.016)

	// Should keep the advanced state since it's the same animation
	if animPtr.FrameIndex != 3 {
		t.Errorf("FrameIndex = %d, want 3 (unchanged)", animPtr.FrameIndex)
	}
}

func TestFacingAnimationStateName(t *testing.T) {
	tests := []struct {
		direction component.FacingDirection
		walking   bool
		want      string
	}{
		{component.FacingDown, false, "idle_down"},
		{component.FacingDown, true, "walk_down"},
		{component.FacingLeft, false, "idle_left"},
		{component.FacingLeft, true, "walk_left"},
		{component.FacingRight, false, "idle_right"},
		{component.FacingRight, true, "walk_right"},
		{component.FacingUp, false, "idle_up"},
		{component.FacingUp, true, "walk_up"},
	}

	for _, tt := range tests {
		f := &component.Facing{Direction: tt.direction, Walking: tt.walking}
		got := facingAnimationStateName(f)
		if got != tt.want {
			t.Errorf("facingAnimationStateName(%v, %v) = %s, want %s", tt.direction, tt.walking, got, tt.want)
		}
	}
}
