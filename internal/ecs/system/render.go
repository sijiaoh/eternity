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

		screenX, screenY := CameraWorldToScreen(s.camera, pos.X, pos.Y)
		DrawSprite(screen, sprite, screenX, screenY)
	})
}

// DrawSprite renders a sprite at the given screen position.
func DrawSprite(screen *ebiten.Image, sprite *component.Sprite, x, y float64) {
	if sprite.Image == nil {
		return
	}

	w, h := sprite.Image.Bounds().Dx(), sprite.Image.Bounds().Dy()
	drawX, drawY := SpriteCalcDrawPosition(sprite, x, y, w, h)

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(drawX, drawY)
	screen.DrawImage(sprite.Image, opts)
}

// SpriteCalcDrawPosition returns draw position adjusted for anchor.
func SpriteCalcDrawPosition(sprite *component.Sprite, x, y float64, imageWidth, imageHeight int) (drawX, drawY float64) {
	offsetX := float64(imageWidth) * sprite.Anchor.X
	offsetY := float64(imageHeight) * sprite.Anchor.Y
	return x - offsetX, y - offsetY
}
