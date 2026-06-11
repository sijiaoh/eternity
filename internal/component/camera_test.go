package component

import (
	"math"
	"testing"

	"ebiten-agent-example/internal/config"
)

func TestCamera_SnapTo(t *testing.T) {
	c := NewCamera(0, 0, 0.1)
	c.SnapTo(5.0, 3.0)

	if c.Position.X != 5.0 || c.Position.Y != 3.0 {
		t.Errorf("SnapTo(5, 3) position = (%v, %v), want (5, 3)", c.Position.X, c.Position.Y)
	}
}

func TestCamera_Update_MovesTowardTarget(t *testing.T) {
	c := NewCamera(0, 0, 0.1) // halfLife = 0.1s
	c.Update(10.0, 10.0, 1.0/60.0)

	// After one frame, camera should move toward target
	if c.Position.X <= 0 || c.Position.Y <= 0 {
		t.Errorf("camera should move toward target, got (%v, %v)", c.Position.X, c.Position.Y)
	}
	if c.Position.X >= 10 || c.Position.Y >= 10 {
		t.Errorf("camera should not reach target in one frame, got (%v, %v)", c.Position.X, c.Position.Y)
	}
}

func TestCamera_Update_NoMovementAtTarget(t *testing.T) {
	c := NewCamera(5.0, 5.0, 0.1)
	c.Update(5.0, 5.0, 1.0/60.0)

	if c.Position.X != 5.0 || c.Position.Y != 5.0 {
		t.Errorf("camera at target should not move, got (%v, %v)", c.Position.X, c.Position.Y)
	}
}

func TestCamera_Update_InstantWithHalfLifeZero(t *testing.T) {
	c := NewCamera(0, 0, 0) // halfLife = 0 means instant
	c.Update(10.0, 5.0, 1.0/60.0)

	if c.Position.X != 10.0 || c.Position.Y != 5.0 {
		t.Errorf("halfLife=0 should reach target instantly, got (%v, %v)", c.Position.X, c.Position.Y)
	}
}

func TestCamera_Update_HalfLifeBehavior(t *testing.T) {
	halfLife := 0.1 // 0.1 seconds to reach halfway
	c := NewCamera(0, 0, halfLife)

	// After exactly halfLife seconds, should be at 50% of the way
	c.Update(10.0, 0, halfLife)

	// Position should be 5.0 (halfway from 0 to 10)
	if math.Abs(c.Position.X-5.0) > 0.0001 {
		t.Errorf("after halfLife seconds, position should be halfway: got %v, want 5.0", c.Position.X)
	}
}

func TestCamera_Update_FrameRateIndependent(t *testing.T) {
	// Two cameras: one updated at 60fps, one at 30fps, for same total time
	halfLife := 0.1
	c60 := NewCamera(0, 0, halfLife)
	c30 := NewCamera(0, 0, halfLife)

	// Simulate 0.5 seconds
	totalTime := 0.5

	// 60fps: 30 frames
	for i := 0; i < 30; i++ {
		c60.Update(10.0, 10.0, 1.0/60.0)
	}

	// 30fps: 15 frames
	for i := 0; i < 15; i++ {
		c30.Update(10.0, 10.0, 1.0/30.0)
	}

	// Both should end up at approximately the same position
	tolerance := 0.0001
	if math.Abs(c60.Position.X-c30.Position.X) > tolerance {
		t.Errorf("X position differs: 60fps=%v, 30fps=%v (diff > %v)", c60.Position.X, c30.Position.X, tolerance)
	}
	if math.Abs(c60.Position.Y-c30.Position.Y) > tolerance {
		t.Errorf("Y position differs: 60fps=%v, 30fps=%v (diff > %v)", c60.Position.Y, c30.Position.Y, tolerance)
	}

	// Verify expected position: after 0.5s with halfLife=0.1s, that's 5 half-lives
	// Remaining distance = 10 * 0.5^5 = 10 * 0.03125 = 0.3125
	// Position = 10 - 0.3125 = 9.6875
	expectedPos := 10.0 * (1 - math.Pow(0.5, totalTime/halfLife))
	if math.Abs(c60.Position.X-expectedPos) > 0.0001 {
		t.Errorf("unexpected position: got %v, expected %v", c60.Position.X, expectedPos)
	}
}

func TestCamera_GetOffset_CameraAtOrigin(t *testing.T) {
	c := NewCamera(0, 0, 0.1) // halfLife doesn't affect offset calculation
	offsetX, offsetY := c.GetOffset()

	// Camera at origin: offset should be -screenCenter (so origin draws at screen center)
	wantX := -float64(config.ScreenWidth) / 2
	wantY := -float64(config.ScreenHeight) / 2

	if offsetX != wantX || offsetY != wantY {
		t.Errorf("GetOffset() = (%v, %v), want (%v, %v)", offsetX, offsetY, wantX, wantY)
	}
}

func TestCamera_WorldToScreen_CameraAtOrigin(t *testing.T) {
	c := NewCamera(0, 0, 0.1) // halfLife doesn't affect coordinate conversion

	tests := []struct {
		name                     string
		worldX, worldY           float64
		wantScreenX, wantScreenY float64
	}{
		{
			"origin at screen center",
			0, 0,
			float64(config.ScreenWidth) / 2, float64(config.ScreenHeight) / 2,
		},
		{
			"1 unit right of origin",
			1, 0,
			float64(config.ScreenWidth)/2 + config.PixelsPerUnit, float64(config.ScreenHeight) / 2,
		},
		{
			"1 unit below origin",
			0, 1,
			float64(config.ScreenWidth) / 2, float64(config.ScreenHeight)/2 + config.PixelsPerUnit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			screenX, screenY := c.WorldToScreen(tt.worldX, tt.worldY)
			if screenX != tt.wantScreenX || screenY != tt.wantScreenY {
				t.Errorf("WorldToScreen(%v, %v) = (%v, %v), want (%v, %v)",
					tt.worldX, tt.worldY, screenX, screenY, tt.wantScreenX, tt.wantScreenY)
			}
		})
	}
}

func TestCamera_WorldToScreen_CameraOffset(t *testing.T) {
	// Camera at (2, 3) units
	c := NewCamera(2, 3, 0.1)

	// World origin (0, 0) should appear shifted left and up from screen center
	screenX, screenY := c.WorldToScreen(0, 0)

	wantX := float64(config.ScreenWidth)/2 - 2*config.PixelsPerUnit
	wantY := float64(config.ScreenHeight)/2 - 3*config.PixelsPerUnit

	if screenX != wantX || screenY != wantY {
		t.Errorf("WorldToScreen(0, 0) with camera at (2,3) = (%v, %v), want (%v, %v)",
			screenX, screenY, wantX, wantY)
	}

	// Camera position should be at screen center
	screenX, screenY = c.WorldToScreen(2, 3)
	if screenX != float64(config.ScreenWidth)/2 || screenY != float64(config.ScreenHeight)/2 {
		t.Errorf("camera position should be at screen center, got (%v, %v)", screenX, screenY)
	}
}
