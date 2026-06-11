package component

import "testing"

func TestPosition_Move(t *testing.T) {
	tests := []struct {
		name         string
		initial      Position
		dx, dy       float64
		wantX, wantY float64
	}{
		{"move positive", Position{X: 0, Y: 0}, 10, 20, 10, 20},
		{"move negative", Position{X: 100, Y: 100}, -30, -50, 70, 50},
		{"no movement", Position{X: 50, Y: 50}, 0, 0, 50, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.initial
			p.Move(tt.dx, tt.dy)
			if p.X != tt.wantX || p.Y != tt.wantY {
				t.Errorf("Move(%v, %v) = (%v, %v), want (%v, %v)", tt.dx, tt.dy, p.X, p.Y, tt.wantX, tt.wantY)
			}
		})
	}
}
