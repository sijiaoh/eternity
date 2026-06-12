//go:build test

package system

import (
	"testing"

	"eternity/internal/config"
)

func TestSpriteSheetCalcScale_NativeSize(t *testing.T) {
	tests := []struct {
		name        string
		frameWidth  int
		sizeInUnits float64
		want        float64
	}{
		{"zero sizeInUnits", 48, 0, 1.0},
		{"negative sizeInUnits", 48, -1.0, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SpriteSheetCalcScale(tt.frameWidth, tt.sizeInUnits)
			if got != tt.want {
				t.Errorf("SpriteSheetCalcScale(%d, %v) = %v, want %v",
					tt.frameWidth, tt.sizeInUnits, got, tt.want)
			}
		})
	}
}

func TestSpriteSheetCalcScale_ScaledSize(t *testing.T) {
	// PixelsPerUnit = 48
	tests := []struct {
		name        string
		frameWidth  int
		sizeInUnits float64
		want        float64
	}{
		{"48px frame at 1 unit", 48, 1.0, 1.0},    // 48 / 48 = 1.0
		{"48px frame at 2 units", 48, 2.0, 2.0},   // 96 / 48 = 2.0
		{"96px frame at 1 unit", 96, 1.0, 0.5},    // 48 / 96 = 0.5
		{"24px frame at 1 unit", 24, 1.0, 2.0},    // 48 / 24 = 2.0
		{"48px frame at 0.5 units", 48, 0.5, 0.5}, // 24 / 48 = 0.5
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SpriteSheetCalcScale(tt.frameWidth, tt.sizeInUnits)
			if got != tt.want {
				t.Errorf("SpriteSheetCalcScale(%d, %v) = %v, want %v",
					tt.frameWidth, tt.sizeInUnits, got, tt.want)
			}
		})
	}
}

func TestSpriteSheetCalcDrawPosition_Center(t *testing.T) {
	// Using center anchor (0.5, 0.5)
	tests := []struct {
		name         string
		x, y, scale  float64
		wantX, wantY float64
	}{
		{"scale 1.0", 100, 100, 1.0, 76, 76},  // 100 - 48*0.5 = 76
		{"scale 2.0", 100, 100, 2.0, 52, 52},  // 100 - 96*0.5 = 52
		{"scale 0.5", 100, 100, 0.5, 88, 88},  // 100 - 24*0.5 = 88
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotY := SpriteSheetCalcDrawPosition(48, 48, 0.5, 0.5, tt.x, tt.y, tt.scale)
			if gotX != tt.wantX || gotY != tt.wantY {
				t.Errorf("SpriteSheetCalcDrawPosition at (%v,%v) scale %v = (%v, %v), want (%v, %v)",
					tt.x, tt.y, tt.scale, gotX, gotY, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestSpriteSheetCalcDrawPosition_TopLeft(t *testing.T) {
	// With top-left anchor (0, 0), offset is always 0
	gotX, gotY := SpriteSheetCalcDrawPosition(48, 48, 0, 0, 100, 100, 2.0)
	if gotX != 100 || gotY != 100 {
		t.Errorf("top-left anchor should have no offset, got (%v, %v)", gotX, gotY)
	}
}

func TestSpriteSheetCalcDrawPosition_BottomRight(t *testing.T) {
	// With bottom-right anchor (1, 1), offset equals scaled size
	// At scale 2.0: offset = 48 * 2 * 1 = 96
	gotX, gotY := SpriteSheetCalcDrawPosition(48, 48, 1, 1, 100, 100, 2.0)
	wantX, wantY := 4.0, 4.0 // 100 - 96 = 4
	if gotX != wantX || gotY != wantY {
		t.Errorf("bottom-right anchor: got (%v, %v), want (%v, %v)", gotX, gotY, wantX, wantY)
	}
}

func TestSpriteSheetCalcDrawPosition_NonSquareFrame(t *testing.T) {
	// Non-square frame (64x32) with center anchor
	// At scale 1.0: offsetX = 64 * 0.5 = 32, offsetY = 32 * 0.5 = 16
	gotX, gotY := SpriteSheetCalcDrawPosition(64, 32, 0.5, 0.5, 100, 100, 1.0)
	wantX, wantY := 68.0, 84.0 // 100 - 32 = 68, 100 - 16 = 84
	if gotX != wantX || gotY != wantY {
		t.Errorf("non-square frame: got (%v, %v), want (%v, %v)", gotX, gotY, wantX, wantY)
	}
}

func TestSpriteSheetCalcScale_MatchesUnitsSystem(t *testing.T) {
	// Verify the scale calculation integrates correctly with config.UnitsToPixels
	frameWidth := 64
	sizeInUnits := 1.5

	scale := SpriteSheetCalcScale(frameWidth, sizeInUnits)
	scaledWidth := float64(frameWidth) * scale
	expectedWidth := config.UnitsToPixels(sizeInUnits)

	if scaledWidth != expectedWidth {
		t.Errorf("scaled width %v should match UnitsToPixels(%v) = %v",
			scaledWidth, sizeInUnits, expectedWidth)
	}
}
