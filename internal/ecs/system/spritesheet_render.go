//go:build !test

package system

import (
	"eternity/internal/component"
	"eternity/internal/ecs"

	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteSheetRenderSystem draws entities with Position, Animation, and SpriteSheet.
type SpriteSheetRenderSystem struct {
	positions    *ecs.Storage[component.Position]
	animations   *ecs.Storage[component.Animation]
	spriteSheets *ecs.Storage[component.SpriteSheet]
	camera       *component.Camera
}

func NewSpriteSheetRenderSystem(
	positions *ecs.Storage[component.Position],
	animations *ecs.Storage[component.Animation],
	spriteSheets *ecs.Storage[component.SpriteSheet],
	camera *component.Camera,
) *SpriteSheetRenderSystem {
	return &SpriteSheetRenderSystem{
		positions:    positions,
		animations:   animations,
		spriteSheets: spriteSheets,
		camera:       camera,
	}
}

func (s *SpriteSheetRenderSystem) Draw(w *ecs.World, screen *ebiten.Image) {
	s.spriteSheets.Each(func(e ecs.Entity, sheet *component.SpriteSheet) {
		if !w.Alive(e) {
			return
		}
		pos, ok := s.positions.Get(e)
		if !ok {
			return
		}
		anim, ok := s.animations.Get(e)
		if !ok {
			return
		}

		screenX, screenY := s.camera.WorldToScreen(pos.X, pos.Y)
		frame := sheet.Frame(anim.Frame())
		if frame == nil {
			return
		}

		drawX, drawY := sheet.CalcDrawPosition(screenX, screenY)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(drawX, drawY)
		screen.DrawImage(frame, opts)
	})
}
