//go:build !test

package system

import (
	"eternity/internal/component"
	"eternity/internal/ecs"
)

// SpriteFlipSystem updates sprite horizontal flip based on facing direction.
// For sprites that only have left-facing frames, this flips them when facing right.
type SpriteFlipSystem struct {
	facings      *ecs.Storage[component.Facing]
	spriteSheets *ecs.Storage[component.SpriteSheet]
}

func NewSpriteFlipSystem(
	facings *ecs.Storage[component.Facing],
	spriteSheets *ecs.Storage[component.SpriteSheet],
) *SpriteFlipSystem {
	return &SpriteFlipSystem{
		facings:      facings,
		spriteSheets: spriteSheets,
	}
}

func (s *SpriteFlipSystem) Update(w *ecs.World, dt float64) {
	s.facings.Each(func(e ecs.Entity, facing *component.Facing) {
		if !w.Alive(e) {
			return
		}
		sheet := s.spriteSheets.GetPtr(e)
		if sheet == nil || !sheet.FlipOnRight {
			return
		}

		sheet.FlipH = facing.Direction == component.FacingRight
	})
}
