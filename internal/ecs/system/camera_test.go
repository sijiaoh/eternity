//go:build !test

package system

import (
	"testing"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

func TestCameraSystem_FollowsTarget(t *testing.T) {
	w := ecs.NewWorld(10)
	targets := ecs.NewStorage[component.CameraTarget](10)
	positions := ecs.NewStorage[component.Position](10)
	camera := component.NewCamera(0, 0, 0) // halfLife=0 for instant follow

	sys := NewCameraSystem(targets, positions, camera)

	e := w.Spawn()
	positions.Set(e, component.Position{X: 10, Y: 20})
	targets.Set(e, component.CameraTarget{})

	sys.Update(w, 0.1)

	if camera.Position.X != 10 || camera.Position.Y != 20 {
		t.Errorf("camera = (%v, %v), want (10, 20)", camera.Position.X, camera.Position.Y)
	}
}

func TestCameraSystem_SkipsDeadTarget(t *testing.T) {
	w := ecs.NewWorld(10)
	targets := ecs.NewStorage[component.CameraTarget](10)
	positions := ecs.NewStorage[component.Position](10)
	camera := component.NewCamera(0, 0, 0)

	sys := NewCameraSystem(targets, positions, camera)

	e := w.Spawn()
	positions.Set(e, component.Position{X: 10, Y: 20})
	targets.Set(e, component.CameraTarget{})

	w.Despawn(e)
	sys.Update(w, 0.1)

	// Camera should not move to despawned target
	if camera.Position.X != 0 || camera.Position.Y != 0 {
		t.Errorf("camera should not follow dead target, got (%v, %v)", camera.Position.X, camera.Position.Y)
	}
}

func TestCameraSystem_NoTarget(t *testing.T) {
	w := ecs.NewWorld(10)
	targets := ecs.NewStorage[component.CameraTarget](10)
	positions := ecs.NewStorage[component.Position](10)
	camera := component.NewCamera(5, 5, 0)

	sys := NewCameraSystem(targets, positions, camera)

	// No target entity
	sys.Update(w, 0.1)

	// Camera should stay at initial position
	if camera.Position.X != 5 || camera.Position.Y != 5 {
		t.Errorf("camera should not move without target, got (%v, %v)", camera.Position.X, camera.Position.Y)
	}
}
