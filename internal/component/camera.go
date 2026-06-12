package component

// Camera holds camera state for world-to-screen transformations.
type Camera struct {
	Position Position
	HalfLife float64 // seconds to reach halfway to target; 0 = instant
}

func NewCamera(x, y, halfLife float64) *Camera {
	return &Camera{
		Position: Position{X: x, Y: y},
		HalfLife: halfLife,
	}
}
