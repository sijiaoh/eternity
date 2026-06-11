package config

// PixelsPerUnit: 1 unit = 48 pixels (matches tile/character size)
const PixelsPerUnit = 48.0

func UnitsToPixels(units float64) float64 {
	return units * PixelsPerUnit
}

func PixelsToUnits(pixels float64) float64 {
	return pixels / PixelsPerUnit
}
