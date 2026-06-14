//go:build !test

package entity

import (
	"eternity/internal/component"
	"eternity/internal/ecs"

	"github.com/hajimehoshi/ebiten/v2"
)

// GoblinComponents holds all storages needed to create a goblin entity.
type GoblinComponents struct {
	Positions    *ecs.Storage[component.Position]
	Velocities   *ecs.Storage[component.Velocity]
	AIFollows    *ecs.Storage[component.AIFollow]
	Facings      *ecs.Storage[component.Facing]
	Animations   *ecs.Storage[component.Animation]
	Directionals *ecs.Storage[component.DirectionalAnimation]
	SpriteSheets *ecs.Storage[component.SpriteSheet]
}

// GoblinFactoryConfig contains parameters for creating a goblin.
type GoblinFactoryConfig struct {
	X, Y        float64
	Speed       float64
	SpriteSheet *ebiten.Image
	FrameWidth  int
	FrameHeight int
	Columns     int
	AnimFPS     float64
	SizeInUnits float64 // Target width in world units; 0 = native size
}

// CreateGoblin spawns a goblin: a generic character plus follow AI. The per-direction
// row layout lives in the DirectionalSheetSpec below (see sprite.source.md for the sheet).
func CreateGoblin(w *ecs.World, c *GoblinComponents, cfg GoblinFactoryConfig) ecs.Entity {
	e := w.Spawn()

	c.Positions.Set(e, component.Position{X: cfg.X, Y: cfg.Y})
	c.Velocities.Set(e, component.Velocity{X: 0, Y: 0})
	// Target is set by AITargetingSystem, not here.
	c.AIFollows.Set(e, component.AIFollow{Speed: cfg.Speed})
	c.Facings.Set(e, component.Facing{Direction: component.FacingDown, Walking: false})

	// Four-directional sheet: 11-frame rows, walk uses the first 8, idle is the row's first frame.
	spec := component.DirectionalSheetSpec{
		Directions: component.DirectionsFour,
		Rows: map[component.FacingDirection]int{
			component.FacingDown:  0,
			component.FacingRight: 11,
			component.FacingUp:    22,
			component.FacingLeft:  33,
		},
		IdleFrames: 1,
		WalkFrames: 8,
		FPS:        cfg.AnimFPS,
	}
	c.Animations.Set(e, *component.NewAnimation(spec.States()))
	c.Directionals.Set(e, *component.NewDirectionalAnimation(spec.Directions))

	c.SpriteSheets.Set(e, component.SpriteSheet{
		Image:       cfg.SpriteSheet,
		FrameWidth:  cfg.FrameWidth,
		FrameHeight: cfg.FrameHeight,
		Columns:     cfg.Columns,
		Anchor:      component.AnchorCenter(),
		SizeInUnits: cfg.SizeInUnits,
	})

	return e
}
