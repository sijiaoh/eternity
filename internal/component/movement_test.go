package component

import (
	"math"
	"testing"
)

func TestMovement_CalcDelta_BasicDirections(t *testing.T) {
	m := NewMovement(100.0) // 100 pixels/second
	dt := 1.0               // 1 second -> expect 100 pixels

	tests := []struct {
		name   string
		dir    MoveDirection
		wantDX float64
		wantDY float64
	}{
		{"no input", MoveDirection{}, 0, 0},
		{"up", MoveDirection{Up: true}, 0, -100},
		{"down", MoveDirection{Down: true}, 0, 100},
		{"left", MoveDirection{Left: true}, -100, 0},
		{"right", MoveDirection{Right: true}, 100, 0},
		{"up+down cancel", MoveDirection{Up: true, Down: true}, 0, 0},
		{"left+right cancel", MoveDirection{Left: true, Right: true}, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dx, dy := m.CalcDelta(tt.dir, dt)
			if dx != tt.wantDX || dy != tt.wantDY {
				t.Errorf("CalcDelta() = (%v, %v), want (%v, %v)", dx, dy, tt.wantDX, tt.wantDY)
			}
		})
	}
}

func TestMovement_DiagonalNormalized(t *testing.T) {
	m := NewMovement(100.0)
	dt := 1.0

	// Diagonal movement should have the same total speed as cardinal movement
	dx, dy := m.CalcDelta(MoveDirection{Up: true, Right: true}, dt)
	actualSpeed := math.Sqrt(dx*dx + dy*dy)

	if math.Abs(actualSpeed-100.0) > 0.001 {
		t.Errorf("diagonal speed = %v, want 100.0 (same as cardinal)", actualSpeed)
	}
}

func TestMovement_DistanceProportionalToTime(t *testing.T) {
	m := NewMovement(60.0) // 60 pixels/second
	dir := MoveDirection{Right: true}

	// Core spec: distance = speed × time
	tests := []struct {
		dt       float64
		wantDist float64
	}{
		{1.0 / 60.0, 1.0}, // 60fps: 1 pixel per frame
		{1.0 / 30.0, 2.0}, // 30fps: 2 pixels per frame
		{1.0, 60.0},       // 1 second: 60 pixels
		{0.5, 30.0},       // 0.5 seconds: 30 pixels
	}

	for _, tt := range tests {
		dx, _ := m.CalcDelta(dir, tt.dt)
		if math.Abs(dx-tt.wantDist) > 0.0001 {
			t.Errorf("dt=%v: got %v, want %v", tt.dt, dx, tt.wantDist)
		}
	}
}

func TestMovement_ZeroDeltaTime(t *testing.T) {
	m := NewMovement(100.0)
	dx, dy := m.CalcDelta(MoveDirection{Up: true, Right: true}, 0)

	if dx != 0 || dy != 0 {
		t.Errorf("zero deltaTime should yield zero movement, got (%v, %v)", dx, dy)
	}
}
