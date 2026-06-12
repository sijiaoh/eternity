//go:build !test

package component

import "testing"

func TestSprite_CalcDrawPosition(t *testing.T) {
	tests := []struct {
		name   string
		anchor Anchor
		x, y   float64
		width  int
		height int
		wantX  float64
		wantY  float64
	}{
		{
			name:   "top-left anchor: no offset",
			anchor: Anchor{X: 0, Y: 0},
			x:      100, y: 100,
			width: 50, height: 40,
			wantX: 100, wantY: 100,
		},
		{
			name:   "center anchor: offset by half size",
			anchor: Anchor{X: 0.5, Y: 0.5},
			x:      100, y: 100,
			width: 50, height: 40,
			wantX: 75, wantY: 80,
		},
		{
			name:   "bottom-right anchor: offset by full size",
			anchor: Anchor{X: 1, Y: 1},
			x:      100, y: 100,
			width: 50, height: 40,
			wantX: 50, wantY: 60,
		},
		{
			name:   "custom anchor",
			anchor: Anchor{X: 0.25, Y: 0.75},
			x:      100, y: 100,
			width: 80, height: 40,
			wantX: 80, wantY: 70,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Sprite{Anchor: tt.anchor}
			gotX, gotY := s.CalcDrawPosition(tt.x, tt.y, tt.width, tt.height)
			if gotX != tt.wantX || gotY != tt.wantY {
				t.Errorf("CalcDrawPosition() = (%v, %v), want (%v, %v)", gotX, gotY, tt.wantX, tt.wantY)
			}
		})
	}
}
