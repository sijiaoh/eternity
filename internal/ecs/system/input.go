//go:build !test

package system

import (
	"eternity/internal/component"
	"eternity/internal/ecs"

	"github.com/hajimehoshi/ebiten/v2"
)

// InputSystem reads keyboard input and updates velocity for InputControlled entities.
type InputSystem struct {
	inputs     *ecs.Storage[component.InputControlled]
	velocities *ecs.Storage[component.Velocity]
}

func NewInputSystem(
	inputs *ecs.Storage[component.InputControlled],
	velocities *ecs.Storage[component.Velocity],
) *InputSystem {
	return &InputSystem{
		inputs:     inputs,
		velocities: velocities,
	}
}

func (s *InputSystem) Update(w *ecs.World, dt float64) {
	dx, dy := readMovementInput()

	s.inputs.Each(func(e ecs.Entity, input *component.InputControlled) {
		if !w.Alive(e) {
			return
		}
		vel := s.velocities.GetPtr(e)
		if vel == nil {
			return
		}
		vel.X = dx * input.Speed
		vel.Y = dy * input.Speed
	})
}

// readMovementInput returns normalized direction based on WASD/arrow keys.
func readMovementInput() (dx, dy float64) {
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		dy -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		dy += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx += 1
	}

	// Normalize diagonal movement
	if dx != 0 && dy != 0 {
		const diagonalFactor = 0.7071067811865476 // 1/sqrt(2)
		dx *= diagonalFactor
		dy *= diagonalFactor
	}

	return dx, dy
}
