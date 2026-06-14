package component

import "testing"

func newTestAnimation() *Animation {
	return NewAnimation([]AnimationState{
		{Name: "idle_down", StartFrame: 0, FrameCount: 1, FPS: 8, Loop: true},
		{Name: "walk_down", StartFrame: 0, FrameCount: 6, FPS: 8, Loop: true},
	})
}

func TestAnimation_NewStartsAtFirstState(t *testing.T) {
	if got := newTestAnimation().CurrentState; got != "idle_down" {
		t.Errorf("initial state = %q, want idle_down (first declared)", got)
	}
}

func TestAnimation_SetStateResetsPlayback(t *testing.T) {
	anim := newTestAnimation()
	anim.FrameIndex = 5
	anim.Elapsed = 0.2
	anim.Finished = true

	anim.SetState("walk_down")

	if anim.CurrentState != "walk_down" {
		t.Errorf("state = %q, want walk_down", anim.CurrentState)
	}
	if anim.FrameIndex != 0 || anim.Elapsed != 0 || anim.Finished {
		t.Errorf("playback not reset: %+v", anim)
	}
}

func TestAnimation_SetStateSameStateKeepsPlayback(t *testing.T) {
	anim := newTestAnimation() // already idle_down
	anim.FrameIndex = 3
	anim.Elapsed = 0.1

	anim.SetState("idle_down")

	if anim.FrameIndex != 3 || anim.Elapsed != 0.1 {
		t.Errorf("same-state switch should not reset playback: %+v", anim)
	}
}

func TestAnimation_SetStatePanicsOnUnknownState(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected panic when switching to an undeclared state")
		}
	}()
	newTestAnimation().SetState("walk_up") // never declared → fail fast
}
