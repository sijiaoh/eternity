//go:build !test

package component

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// SpriteSheet holds sprite sheet data for animated entities.
type SpriteSheet struct {
	Image       *ebiten.Image
	FrameWidth  int
	FrameHeight int
	Columns     int
	Anchor      Anchor
}

func NewSpriteSheet(img *ebiten.Image, frameWidth, frameHeight, columns int) *SpriteSheet {
	return &SpriteSheet{
		Image:       img,
		FrameWidth:  frameWidth,
		FrameHeight: frameHeight,
		Columns:     columns,
		Anchor:      AnchorCenter(),
	}
}

// Frame returns the sub-image for the given frame index.
func (s *SpriteSheet) Frame(index int) *ebiten.Image {
	if s.Image == nil {
		return nil
	}

	col := index % s.Columns
	row := index / s.Columns

	x := col * s.FrameWidth
	y := row * s.FrameHeight

	rect := image.Rect(x, y, x+s.FrameWidth, y+s.FrameHeight)
	return s.Image.SubImage(rect).(*ebiten.Image)
}

// CalcDrawPosition returns draw position adjusted for anchor.
func (s *SpriteSheet) CalcDrawPosition(x, y float64) (drawX, drawY float64) {
	offsetX := float64(s.FrameWidth) * s.Anchor.X
	offsetY := float64(s.FrameHeight) * s.Anchor.Y
	return x - offsetX, y - offsetY
}
