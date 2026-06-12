//go:build !test

package system

import (
	"image"

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

		screenX, screenY := CameraWorldToScreen(s.camera, pos.X, pos.Y)
		frame := SpriteSheetFrame(sheet, anim.Frame())
		if frame == nil {
			return
		}

		drawX, drawY := SpriteSheetCalcDrawPosition(sheet, screenX, screenY)
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(drawX, drawY)
		screen.DrawImage(frame, opts)
	})
}

// SpriteSheetFrame returns the sub-image for the given frame index.
func SpriteSheetFrame(sheet *component.SpriteSheet, index int) *ebiten.Image {
	if sheet.Image == nil {
		return nil
	}

	col := index % sheet.Columns
	row := index / sheet.Columns

	x := col * sheet.FrameWidth
	y := row * sheet.FrameHeight

	rect := image.Rect(x, y, x+sheet.FrameWidth, y+sheet.FrameHeight)
	return sheet.Image.SubImage(rect).(*ebiten.Image)
}

// SpriteSheetCalcDrawPosition returns draw position adjusted for anchor.
func SpriteSheetCalcDrawPosition(sheet *component.SpriteSheet, x, y float64) (drawX, drawY float64) {
	offsetX := float64(sheet.FrameWidth) * sheet.Anchor.X
	offsetY := float64(sheet.FrameHeight) * sheet.Anchor.Y
	return x - offsetX, y - offsetY
}
