package component

import (
	"math"

	"ebiten-agent-example/internal/config"
)

// Camera follows a target position and provides world-to-screen conversion.
type Camera struct {
	Position Position
	halfLife float64 // seconds to reach halfway to target; 0 = instant
}

func NewCamera(x, y, halfLife float64) *Camera {
	return &Camera{
		Position: Position{X: x, Y: y},
		halfLife: halfLife,
	}
}

// SnapTo sets the camera position directly (no smoothing).
func (c *Camera) SnapTo(x, y float64) {
	c.Position.X = x
	c.Position.Y = y
}

// Update smoothly moves camera toward target position.
// Uses exponential decay for frame-rate independent behavior.
func (c *Camera) Update(targetX, targetY, deltaTime float64) {
	if c.halfLife <= 0 {
		c.Position.X = targetX
		c.Position.Y = targetY
		return
	}
	// Exponential decay: after halfLife seconds, distance is halved
	decay := math.Pow(0.5, deltaTime/c.halfLife)
	c.Position.X = targetX + (c.Position.X-targetX)*decay
	c.Position.Y = targetY + (c.Position.Y-targetY)*decay
}

// GetOffset returns the rendering offset in pixels.
// Objects should be drawn at (worldPixelX - offsetX, worldPixelY - offsetY).
func (c *Camera) GetOffset() (offsetX, offsetY float64) {
	offsetX = config.UnitsToPixels(c.Position.X) - float64(config.ScreenWidth)/2
	offsetY = config.UnitsToPixels(c.Position.Y) - float64(config.ScreenHeight)/2
	return
}

// WorldToScreen converts a world position to screen pixel coordinates.
func (c *Camera) WorldToScreen(worldX, worldY float64) (screenX, screenY float64) {
	offsetX, offsetY := c.GetOffset()
	screenX = config.UnitsToPixels(worldX) - offsetX
	screenY = config.UnitsToPixels(worldY) - offsetY
	return
}
