package system

import "eternity/internal/config"

// SpriteSheetCalcScale returns the scale factor to render at the specified world size.
// Returns 1.0 if sizeInUnits <= 0 (native size).
func SpriteSheetCalcScale(frameWidth int, sizeInUnits float64) float64 {
	if sizeInUnits <= 0 {
		return 1.0
	}
	return config.UnitsToPixels(sizeInUnits) / float64(frameWidth)
}

// SpriteSheetCalcDrawPosition returns draw position adjusted for anchor and scale.
// anchorX and anchorY are in 0-1 range (e.g., 0.5 for center).
func SpriteSheetCalcDrawPosition(frameWidth, frameHeight int, anchorX, anchorY, x, y, scale float64) (drawX, drawY float64) {
	scaledW := float64(frameWidth) * scale
	scaledH := float64(frameHeight) * scale
	offsetX := scaledW * anchorX
	offsetY := scaledH * anchorY
	return x - offsetX, y - offsetY
}
