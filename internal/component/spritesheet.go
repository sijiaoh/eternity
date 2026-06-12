//go:build !test

package component

import "github.com/hajimehoshi/ebiten/v2"

// SpriteSheet holds sprite sheet data for animated entities.
type SpriteSheet struct {
	Image       *ebiten.Image
	FrameWidth  int
	FrameHeight int
	Columns     int
	Anchor      Anchor
	FlipH       bool    // Horizontal flip for rendering (set by SpriteFlipSystem)
	FlipOnRight bool    // If true, SpriteFlipSystem sets FlipH=true when facing right
	SizeInUnits float64 // Target width in world units; 0 = native size
}

func NewSpriteSheet(img *ebiten.Image, frameWidth, frameHeight, columns int) *SpriteSheet {
	return &SpriteSheet{
		Image:       img,
		FrameWidth:  frameWidth,
		FrameHeight: frameHeight,
		Columns:     columns,
		Anchor:      AnchorCenter(),
		FlipH:       false,
	}
}
