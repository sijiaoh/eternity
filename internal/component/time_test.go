package component

import (
	"math"
	"testing"
)

func TestNewClock(t *testing.T) {
	c := NewClock()

	if c.Scale() != 1.0 {
		t.Errorf("Scale() = %v, want 1.0", c.Scale())
	}
	if c.DeltaTime() != 0 {
		t.Errorf("DeltaTime() = %v, want 0", c.DeltaTime())
	}
	if c.TotalTime() != 0 {
		t.Errorf("TotalTime() = %v, want 0", c.TotalTime())
	}
}

func TestClock_Update(t *testing.T) {
	c := NewClock()

	c.Update(0.016) // ~60fps frame
	if math.Abs(c.DeltaTime()-0.016) > 0.0001 {
		t.Errorf("DeltaTime() = %v, want 0.016", c.DeltaTime())
	}
	if math.Abs(c.TotalTime()-0.016) > 0.0001 {
		t.Errorf("TotalTime() = %v, want 0.016", c.TotalTime())
	}

	c.Update(0.016)
	if math.Abs(c.TotalTime()-0.032) > 0.0001 {
		t.Errorf("TotalTime() = %v, want 0.032", c.TotalTime())
	}
}

func TestClock_Scale(t *testing.T) {
	tests := []struct {
		name      string
		scale     float64
		rawDelta  float64
		wantDelta float64
	}{
		{"normal speed", 1.0, 0.016, 0.016},
		{"double speed", 2.0, 0.016, 0.032},
		{"half speed", 0.5, 0.016, 0.008},
		{"paused", 0.0, 0.016, 0.0},
		{"triple speed", 3.0, 0.016, 0.048},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClock()
			c.SetScale(tt.scale)
			c.Update(tt.rawDelta)

			if math.Abs(c.DeltaTime()-tt.wantDelta) > 0.0001 {
				t.Errorf("DeltaTime() = %v, want %v", c.DeltaTime(), tt.wantDelta)
			}
		})
	}
}

func TestClock_SetScaleNegative(t *testing.T) {
	c := NewClock()
	c.SetScale(-1.0)

	if c.Scale() != 0 {
		t.Errorf("negative scale should be clamped to 0, got %v", c.Scale())
	}
}

func TestClock_PauseResume(t *testing.T) {
	c := NewClock()

	c.Pause()
	if !c.IsPaused() {
		t.Error("IsPaused() = false after Pause()")
	}
	if c.Scale() != 0 {
		t.Errorf("Scale() = %v after Pause(), want 0", c.Scale())
	}

	c.Update(0.016)
	if c.DeltaTime() != 0 {
		t.Errorf("DeltaTime() = %v while paused, want 0", c.DeltaTime())
	}

	c.Resume()
	if c.IsPaused() {
		t.Error("IsPaused() = true after Resume()")
	}
	if c.Scale() != 1.0 {
		t.Errorf("Scale() = %v after Resume(), want 1.0", c.Scale())
	}
}

func TestClock_PauseResumePreservesScale(t *testing.T) {
	c := NewClock()
	c.SetScale(0.5) // slow motion

	c.Pause()
	if !c.IsPaused() {
		t.Error("IsPaused() = false after Pause()")
	}

	c.Resume()
	if c.Scale() != 0.5 {
		t.Errorf("Scale() = %v after Resume(), want 0.5 (original scale)", c.Scale())
	}
}

func TestNewChildClock(t *testing.T) {
	parent := NewClock()
	child := NewChildClock(parent)

	if child.Scale() != 1.0 {
		t.Errorf("child Scale() = %v, want 1.0", child.Scale())
	}
}

func TestChildClock_InheritsDeltaTime(t *testing.T) {
	parent := NewClock()
	child := NewChildClock(parent)

	parent.Update(0.016)
	child.Tick() // rawDelta ignored for child clocks

	if math.Abs(child.DeltaTime()-0.016) > 0.0001 {
		t.Errorf("child DeltaTime() = %v, want 0.016", child.DeltaTime())
	}
}

func TestChildClock_ScaleStacks(t *testing.T) {
	parent := NewClock()
	child := NewChildClock(parent)

	parent.SetScale(2.0) // 2x speed
	child.SetScale(0.5)  // half of parent

	parent.Update(0.016)
	child.Tick()

	// parent delta = 0.016 * 2 = 0.032
	// child delta = 0.032 * 0.5 = 0.016
	if math.Abs(child.DeltaTime()-0.016) > 0.0001 {
		t.Errorf("child DeltaTime() = %v, want 0.016", child.DeltaTime())
	}
}

func TestChildClock_ParentPausedStopsChild(t *testing.T) {
	parent := NewClock()
	child := NewChildClock(parent)

	parent.Pause()
	parent.Update(0.016)
	child.Tick()

	if child.DeltaTime() != 0 {
		t.Errorf("child DeltaTime() = %v when parent paused, want 0", child.DeltaTime())
	}
}

func TestChildClock_IndependentPause(t *testing.T) {
	parent := NewClock()
	child := NewChildClock(parent)

	child.Pause()
	parent.Update(0.016)
	child.Tick()

	// Parent runs normally
	if math.Abs(parent.DeltaTime()-0.016) > 0.0001 {
		t.Errorf("parent DeltaTime() = %v, want 0.016", parent.DeltaTime())
	}
	// Child is paused independently
	if child.DeltaTime() != 0 {
		t.Errorf("child DeltaTime() = %v when child paused, want 0", child.DeltaTime())
	}
}

func TestNestedClocks_SiblingBranches(t *testing.T) {
	root := NewClock()
	game := NewChildClock(root)
	ui := NewChildClock(root)

	root.SetScale(1.0)
	game.SetScale(0.5) // slow motion
	ui.SetScale(1.0)   // UI runs at normal speed

	root.Update(0.016)
	game.Tick()
	ui.Tick()

	// game delta = 0.016 * 0.5 = 0.008
	if math.Abs(game.DeltaTime()-0.008) > 0.0001 {
		t.Errorf("game DeltaTime() = %v, want 0.008", game.DeltaTime())
	}
	// ui delta = 0.016 * 1.0 = 0.016
	if math.Abs(ui.DeltaTime()-0.016) > 0.0001 {
		t.Errorf("ui DeltaTime() = %v, want 0.016", ui.DeltaTime())
	}
}

func TestNestedClocks_ThreeLevelsDeep(t *testing.T) {
	root := NewClock()
	world := NewChildClock(root)
	entity := NewChildClock(world)

	root.SetScale(1.0)
	world.SetScale(2.0)  // world runs at 2x
	entity.SetScale(0.5) // entity runs at half of world

	root.Update(0.016)
	world.Tick()
	entity.Tick()

	// root delta = 0.016
	// world delta = 0.016 * 2.0 = 0.032
	// entity delta = 0.032 * 0.5 = 0.016
	if math.Abs(entity.DeltaTime()-0.016) > 0.0001 {
		t.Errorf("entity DeltaTime() = %v, want 0.016", entity.DeltaTime())
	}
}

func TestClock_TotalTimeAccumulates(t *testing.T) {
	c := NewClock()
	c.SetScale(2.0)

	for i := 0; i < 10; i++ {
		c.Update(0.016)
	}

	// 10 frames * 0.016 * 2.0 = 0.32
	if math.Abs(c.TotalTime()-0.32) > 0.0001 {
		t.Errorf("TotalTime() = %v, want 0.32", c.TotalTime())
	}
}

func TestChildClock_TotalTimeAccumulates(t *testing.T) {
	parent := NewClock()
	child := NewChildClock(parent)
	child.SetScale(0.5)

	for i := 0; i < 10; i++ {
		parent.Update(0.016)
		child.Tick()
	}

	// child: 10 frames * 0.016 * 0.5 = 0.08
	if math.Abs(child.TotalTime()-0.08) > 0.0001 {
		t.Errorf("child TotalTime() = %v, want 0.08", child.TotalTime())
	}
}
