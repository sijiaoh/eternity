package component

import (
	"testing"
)

func TestNewAnimation(t *testing.T) {
	states := []AnimationState{
		{Name: "idle", StartFrame: 0, FrameCount: 6, FPS: 10, Loop: true},
		{Name: "walk", StartFrame: 6, FrameCount: 6, FPS: 12, Loop: true},
	}
	a := NewAnimation(states)

	if a.State() != "idle" {
		t.Errorf("State() = %q, want %q", a.State(), "idle")
	}
	if a.Frame() != 0 {
		t.Errorf("Frame() = %d, want 0", a.Frame())
	}
}

func TestNewAnimation_Empty(t *testing.T) {
	a := NewAnimation([]AnimationState{})

	if a.State() != "" {
		t.Errorf("State() = %q, want empty", a.State())
	}
	if a.Frame() != 0 {
		t.Errorf("Frame() = %d, want 0", a.Frame())
	}
}

func TestAnimation_Update_AdvancesFrame(t *testing.T) {
	a := NewAnimation([]AnimationState{
		{Name: "walk", StartFrame: 0, FrameCount: 4, FPS: 10, Loop: true},
	})

	// At 10 FPS, each frame lasts 0.1 seconds
	// After 0.1s, we should be on frame 1
	a.Update(0.1)
	if a.Frame() != 1 {
		t.Errorf("Frame() = %d, want 1 after 0.1s", a.Frame())
	}

	// After another 0.1s, we should be on frame 2
	a.Update(0.1)
	if a.Frame() != 2 {
		t.Errorf("Frame() = %d, want 2 after 0.2s", a.Frame())
	}
}

func TestAnimation_Update_RespectsFPS(t *testing.T) {
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
			a := NewAnimation([]AnimationState{
				{Name: "test", StartFrame: 0, FrameCount: 10, FPS: tt.fps, Loop: true},
			})
			a.Update(tt.deltaTime)
			if a.Frame() != tt.wantFrame {
				t.Errorf("Frame() = %d, want %d", a.Frame(), tt.wantFrame)
			}
		})
	}
}

func TestAnimation_Loop(t *testing.T) {
	a := NewAnimation([]AnimationState{
		{Name: "walk", StartFrame: 0, FrameCount: 3, FPS: 10, Loop: true},
	})

	// At 10 FPS with 3 frames, loop every 0.3s
	// Use 0.31s to ensure we pass the 0.3s threshold and trigger loop
	a.Update(0.31)
	if a.Frame() != 0 {
		t.Errorf("Frame() = %d, want 0 (looped)", a.Frame())
	}
	if a.IsFinished() {
		t.Error("looping animation should not be finished")
	}

	// After 0.1s more, we should be on frame 1
	a.Update(0.1)
	if a.Frame() != 1 {
		t.Errorf("Frame() = %d, want 1", a.Frame())
	}
}

func TestAnimation_NoLoop(t *testing.T) {
	a := NewAnimation([]AnimationState{
		{Name: "attack", StartFrame: 0, FrameCount: 3, FPS: 10, Loop: false},
	})

	// After enough time to go past all frames
	a.Update(0.5)
	if a.Frame() != 2 {
		t.Errorf("Frame() = %d, want 2 (last frame)", a.Frame())
	}
	if !a.IsFinished() {
		t.Error("non-looping animation should be finished")
	}

	// Further updates should not change frame
	a.Update(0.5)
	if a.Frame() != 2 {
		t.Errorf("Frame() = %d, want 2 after finished", a.Frame())
	}
}

func TestAnimation_SetState(t *testing.T) {
	a := NewAnimation([]AnimationState{
		{Name: "idle", StartFrame: 0, FrameCount: 6, FPS: 10, Loop: true},
		{Name: "walk", StartFrame: 6, FrameCount: 6, FPS: 12, Loop: true},
	})

	// Advance idle animation
	a.Update(0.15) // frame 1 in idle

	a.SetState("walk")
	if a.State() != "walk" {
		t.Errorf("State() = %q, want %q", a.State(), "walk")
	}
	if a.Frame() != 6 {
		t.Errorf("Frame() = %d, want 6 (walk start)", a.Frame())
	}

	// Advance walk animation
	a.Update(1.0 / 12.0)
	if a.Frame() != 7 {
		t.Errorf("Frame() = %d, want 7", a.Frame())
	}
}

func TestAnimation_SetState_SameStateNoReset(t *testing.T) {
	a := NewAnimation([]AnimationState{
		{Name: "walk", StartFrame: 0, FrameCount: 6, FPS: 10, Loop: true},
	})

	a.Update(0.15) // frame 1
	a.SetState("walk")

	// Setting same state should not reset
	if a.Frame() != 1 {
		t.Errorf("Frame() = %d, want 1 (unchanged)", a.Frame())
	}
}

func TestAnimation_SetState_UnknownState(t *testing.T) {
	a := NewAnimation([]AnimationState{
		{Name: "idle", StartFrame: 0, FrameCount: 6, FPS: 10, Loop: true},
	})

	a.Update(0.1)
	a.SetState("unknown")

	if a.State() != "idle" {
		t.Errorf("State() = %q, want %q (unchanged)", a.State(), "idle")
	}
}

func TestAnimation_Reset(t *testing.T) {
	a := NewAnimation([]AnimationState{
		{Name: "attack", StartFrame: 0, FrameCount: 3, FPS: 10, Loop: false},
	})

	a.Update(0.5)
	if !a.IsFinished() {
		t.Error("should be finished before reset")
	}

	a.Reset()
	if a.Frame() != 0 {
		t.Errorf("Frame() = %d, want 0 after reset", a.Frame())
	}
	if a.IsFinished() {
		t.Error("should not be finished after reset")
	}
}

func TestAnimation_SingleFrame(t *testing.T) {
	a := NewAnimation([]AnimationState{
		{Name: "static", StartFrame: 5, FrameCount: 1, FPS: 10, Loop: true},
	})

	a.Update(1.0)
	if a.Frame() != 5 {
		t.Errorf("Frame() = %d, want 5", a.Frame())
	}
}

func TestAnimation_StartFrame_Offset(t *testing.T) {
	a := NewAnimation([]AnimationState{
		{Name: "row2", StartFrame: 6, FrameCount: 6, FPS: 10, Loop: true},
	})

	if a.Frame() != 6 {
		t.Errorf("initial Frame() = %d, want 6", a.Frame())
	}

	a.Update(0.2) // frame 2 in the animation
	if a.Frame() != 8 {
		t.Errorf("Frame() = %d, want 8 (6+2)", a.Frame())
	}
}

func TestAnimation_ZeroDeltaTime(t *testing.T) {
	a := NewAnimation([]AnimationState{
		{Name: "walk", StartFrame: 0, FrameCount: 6, FPS: 10, Loop: true},
	})

	a.Update(0)
	if a.Frame() != 0 {
		t.Errorf("Frame() = %d, want 0", a.Frame())
	}
}

func TestAnimation_NegativeDeltaTime(t *testing.T) {
	a := NewAnimation([]AnimationState{
		{Name: "walk", StartFrame: 0, FrameCount: 6, FPS: 10, Loop: true},
	})

	a.Update(0.1) // advance to frame 1
	a.Update(-0.5) // negative time should be ignored

	if a.Frame() != 1 {
		t.Errorf("Frame() = %d, want 1 (unchanged)", a.Frame())
	}
}

func TestAnimation_ZeroFPS(t *testing.T) {
	a := NewAnimation([]AnimationState{
		{Name: "broken", StartFrame: 0, FrameCount: 6, FPS: 0, Loop: true},
	})

	// Should not panic
	a.Update(1.0)
	if a.Frame() != 0 {
		t.Errorf("Frame() = %d, want 0", a.Frame())
	}
}

func TestAnimation_MultipleSmallUpdates(t *testing.T) {
	// Test that many small updates accumulate correctly
	a := NewAnimation([]AnimationState{
		{Name: "walk", StartFrame: 0, FrameCount: 6, FPS: 10, Loop: true},
	})

	// 60 updates of 1/60s each = 1 second
	// At 10 FPS, this is 10 frame advances (with floating point variance)
	// 10 % 6 = 4, but due to accumulated floating point error, may be 9 advances = 3
	for i := 0; i < 60; i++ {
		a.Update(1.0 / 60.0)
	}

	// Floating point accumulation may cause 9 or 10 frame advances
	// 9 % 6 = 3, 10 % 6 = 4
	got := a.Frame()
	if got != 3 && got != 4 {
		t.Errorf("Frame() = %d, want 3 or 4", got)
	}
}
