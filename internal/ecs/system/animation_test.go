//go:build !test

package system

import (
	"testing"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

func TestAnimationSystem_Update(t *testing.T) {
	w := ecs.NewWorld(10)
	animations := ecs.NewStorage[component.Animation](10)

	sys := NewAnimationSystem(animations)

	e := w.Spawn()
	anim := *component.NewAnimation([]component.AnimationState{
		{Name: "walk", StartFrame: 0, FrameCount: 4, FPS: 10, Loop: true},
	})
	animations.Set(e, anim)

	// At 10 FPS, each frame = 0.1s. After 0.15s, should be on frame 1
	sys.Update(w, 0.15)

	animPtr := animations.GetPtr(e)
	if animPtr.Frame() != 1 {
		t.Errorf("frame = %d, want 1", animPtr.Frame())
	}
}

func TestAnimationSystem_SkipsDeadEntities(t *testing.T) {
	w := ecs.NewWorld(10)
	animations := ecs.NewStorage[component.Animation](10)

	sys := NewAnimationSystem(animations)

	e := w.Spawn()
	anim := *component.NewAnimation([]component.AnimationState{
		{Name: "walk", StartFrame: 0, FrameCount: 4, FPS: 10, Loop: true},
	})
	animations.Set(e, anim)

	w.Despawn(e)
	sys.Update(w, 1.0)

	// Should not advance for despawned entity
	animPtr := animations.GetPtr(e)
	if animPtr.Frame() != 0 {
		t.Errorf("despawned entity animation should not advance, frame = %d", animPtr.Frame())
	}
}
