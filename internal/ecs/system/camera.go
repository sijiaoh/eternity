package system

import (
	"eternity/internal/component"
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
		s.camera.Update(targetPos.X, targetPos.Y, dt)
	}
}
