//go:build !test

package system

import (
	"eternity/internal/component"
	"eternity/internal/ecs"

	"github.com/hajimehoshi/ebiten/v2"
)

// RenderSystem draws all entities with Position and Sprite components.
type RenderSystem struct {
	positions *ecs.Storage[component.Position]
	sprites   *ecs.Storage[component.Sprite]
	camera    *component.Camera
}

func NewRenderSystem(
	positions *ecs.Storage[component.Position],
	sprites *ecs.Storage[component.Sprite],
	camera *component.Camera,
) *RenderSystem {
	return &RenderSystem{
		positions: positions,
		sprites:   sprites,
		camera:    camera,
	}
}

func (s *RenderSystem) Draw(w *ecs.World, screen *ebiten.Image) {
	s.sprites.Each(func(e ecs.Entity, sprite *component.Sprite) {
		if !w.Alive(e) {
			return
		}
		pos, ok := s.positions.Get(e)
		if !ok {
			return
		}

		screenX, screenY := s.camera.WorldToScreen(pos.X, pos.Y)
		sprite.Draw(screen, screenX, screenY)
	})
}
