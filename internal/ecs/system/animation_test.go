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

func TestAnimationSystem_RespectsFPS(t *testing.T) {
	tests := []struct {
		name      string
		fps       float64
		deltaTime float64
		wantFrame int
	}{
		{"10fps after 0.05s", 10, 0.05, 0}, // not enough time
		{"10fps after 0.1s", 10, 0.1, 1},   // exactly one frame
		{"10fps after 0.25s", 10, 0.25, 2}, // two frames + partial
		{"60fps after 1/60s", 60, 1.0 / 60.0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := ecs.NewWorld(10)
			animations := ecs.NewStorage[component.Animation](10)
			sys := NewAnimationSystem(animations)

			e := w.Spawn()
			anim := *component.NewAnimation([]component.AnimationState{
				{Name: "test", StartFrame: 0, FrameCount: 10, FPS: tt.fps, Loop: true},
			})
			animations.Set(e, anim)

			sys.Update(w, tt.deltaTime)

			animPtr := animations.GetPtr(e)
			if animPtr.Frame() != tt.wantFrame {
				t.Errorf("Frame() = %d, want %d", animPtr.Frame(), tt.wantFrame)
			}
		})
	}
}

func TestAnimationSystem_Loop(t *testing.T) {
	w := ecs.NewWorld(10)
	animations := ecs.NewStorage[component.Animation](10)
	sys := NewAnimationSystem(animations)

	e := w.Spawn()
	anim := *component.NewAnimation([]component.AnimationState{
		{Name: "walk", StartFrame: 0, FrameCount: 3, FPS: 10, Loop: true},
	})
	animations.Set(e, anim)

	// At 10 FPS with 3 frames, loop every 0.3s
	// Use 0.31s to ensure we pass the 0.3s threshold and trigger loop
	sys.Update(w, 0.31)

	animPtr := animations.GetPtr(e)
	if animPtr.Frame() != 0 {
		t.Errorf("Frame() = %d, want 0 (looped)", animPtr.Frame())
	}
	if animPtr.Finished {
		t.Error("looping animation should not be finished")
	}
}

func TestAnimationSystem_NoLoop(t *testing.T) {
	w := ecs.NewWorld(10)
	animations := ecs.NewStorage[component.Animation](10)
	sys := NewAnimationSystem(animations)

	e := w.Spawn()
	anim := *component.NewAnimation([]component.AnimationState{
		{Name: "attack", StartFrame: 0, FrameCount: 3, FPS: 10, Loop: false},
	})
	animations.Set(e, anim)

	// After enough time to go past all frames
	sys.Update(w, 0.5)

	animPtr := animations.GetPtr(e)
	if animPtr.Frame() != 2 {
		t.Errorf("Frame() = %d, want 2 (last frame)", animPtr.Frame())
	}
	if !animPtr.Finished {
		t.Error("non-looping animation should be finished")
	}

	// Further updates should not change frame
	sys.Update(w, 0.5)
	if animPtr.Frame() != 2 {
		t.Errorf("Frame() = %d, want 2 after finished", animPtr.Frame())
	}
}

func TestAnimationSystem_SingleFrame(t *testing.T) {
	w := ecs.NewWorld(10)
	animations := ecs.NewStorage[component.Animation](10)
	sys := NewAnimationSystem(animations)

	e := w.Spawn()
	anim := *component.NewAnimation([]component.AnimationState{
		{Name: "static", StartFrame: 5, FrameCount: 1, FPS: 10, Loop: true},
	})
	animations.Set(e, anim)

	sys.Update(w, 1.0)

	animPtr := animations.GetPtr(e)
	if animPtr.Frame() != 5 {
		t.Errorf("Frame() = %d, want 5", animPtr.Frame())
	}
}

func TestAnimationSystem_ZeroDeltaTime(t *testing.T) {
	w := ecs.NewWorld(10)
	animations := ecs.NewStorage[component.Animation](10)
	sys := NewAnimationSystem(animations)

	e := w.Spawn()
	anim := *component.NewAnimation([]component.AnimationState{
		{Name: "walk", StartFrame: 0, FrameCount: 6, FPS: 10, Loop: true},
	})
	animations.Set(e, anim)

	sys.Update(w, 0)

	animPtr := animations.GetPtr(e)
	if animPtr.Frame() != 0 {
		t.Errorf("Frame() = %d, want 0", animPtr.Frame())
	}
}
