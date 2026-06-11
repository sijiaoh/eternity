package component

import (
	"math"
	"testing"
)

func TestMovement_CalcDelta(t *testing.T) {
	m := NewMovement(10.0)

	tests := []struct {
		name       string
		dir        MoveDirection
		wantDX     float64
		wantDY     float64
		approxDiag bool
	}{
		{"no input", MoveDirection{}, 0, 0, false},
		{"up only", MoveDirection{Up: true}, 0, -10, false},
		{"down only", MoveDirection{Down: true}, 0, 10, false},
		{"left only", MoveDirection{Left: true}, -10, 0, false},
		{"right only", MoveDirection{Right: true}, 10, 0, false},
		{"up+right diagonal", MoveDirection{Up: true, Right: true}, 7.07, -7.07, true},
		{"down+left diagonal", MoveDirection{Down: true, Left: true}, -7.07, 7.07, true},
		{"up+down cancel", MoveDirection{Up: true, Down: true}, 0, 0, false},
		{"left+right cancel", MoveDirection{Left: true, Right: true}, 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dx, dy := m.CalcDelta(tt.dir)
			if tt.approxDiag {
				if math.Abs(dx-tt.wantDX) > 0.01 || math.Abs(dy-tt.wantDY) > 0.01 {
					t.Errorf("CalcDelta() = (%v, %v), want approx (%v, %v)", dx, dy, tt.wantDX, tt.wantDY)
				}
			} else {
				if dx != tt.wantDX || dy != tt.wantDY {
					t.Errorf("CalcDelta() = (%v, %v), want (%v, %v)", dx, dy, tt.wantDX, tt.wantDY)
				}
			}
		})
	}
}

func TestMovement_DiagonalNormalized(t *testing.T) {
	m := NewMovement(10.0)

	dx, dy := m.CalcDelta(MoveDirection{Up: true, Right: true})
	speed := math.Sqrt(dx*dx + dy*dy)

	if math.Abs(speed-10.0) > 0.001 {
		t.Errorf("diagonal speed = %v, want 10.0", speed)
	}
}
