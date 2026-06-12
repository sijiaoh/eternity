package system

import (
	"math"

	"eternity/internal/component"
	"eternity/internal/config"
	"eternity/internal/ecs"
)

// CameraSystem updates the camera to follow the entity with CameraTarget component.
type CameraSystem struct {
	targets   *ecs.Storage[component.CameraTarget]
	positions *ecs.Storage[component.Position]
	camera    *component.Camera
}

func NewCameraSystem(
	targets *ecs.Storage[component.CameraTarget],
	positions *ecs.Storage[component.Position],
	camera *component.Camera,
) *CameraSystem {
	return &CameraSystem{
		targets:   targets,
		positions: positions,
		camera:    camera,
	}
}

func (s *CameraSystem) Update(w *ecs.World, dt float64) {
	// Find the first alive entity with CameraTarget
	var targetPos *component.Position
	s.targets.Each(func(e ecs.Entity, _ *component.CameraTarget) {
		if targetPos != nil {
			return // Already found a target
		}
		if !w.Alive(e) {
			return
		}
		pos := s.positions.GetPtr(e)
		if pos != nil {
			targetPos = pos
		}
	})

	if targetPos != nil {
		UpdateCamera(s.camera, targetPos.X, targetPos.Y, dt)
	}
}

// UpdateCamera smoothly moves camera toward target position.
// Uses exponential decay for frame-rate independent behavior.
func UpdateCamera(c *component.Camera, targetX, targetY, dt float64) {
	if c.HalfLife <= 0 {
		c.Position.X = targetX
		c.Position.Y = targetY
		return
	}
	// Exponential decay: after HalfLife seconds, distance is halved
	decay := math.Pow(0.5, dt/c.HalfLife)
	c.Position.X = targetX + (c.Position.X-targetX)*decay
	c.Position.Y = targetY + (c.Position.Y-targetY)*decay
}

// CameraGetOffset returns the rendering offset in pixels.
// Objects should be drawn at (worldPixelX - offsetX, worldPixelY - offsetY).
func CameraGetOffset(c *component.Camera) (offsetX, offsetY float64) {
	offsetX = config.UnitsToPixels(c.Position.X) - float64(config.ScreenWidth)/2
	offsetY = config.UnitsToPixels(c.Position.Y) - float64(config.ScreenHeight)/2
	return
}

// CameraWorldToScreen converts a world position to screen pixel coordinates.
func CameraWorldToScreen(c *component.Camera, worldX, worldY float64) (screenX, screenY float64) {
	offsetX, offsetY := CameraGetOffset(c)
	screenX = config.UnitsToPixels(worldX) - offsetX
	screenY = config.UnitsToPixels(worldY) - offsetY
	return
}

// SnapCamera sets the camera position directly (no smoothing).
func SnapCamera(c *component.Camera, x, y float64) {
	c.Position.X = x
	c.Position.Y = y
}
