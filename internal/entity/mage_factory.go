//go:build !test

package entity

import (
	"eternity/internal/component"
	"eternity/internal/ecs"

	"github.com/hajimehoshi/ebiten/v2"
)

// MageComponents holds all storages needed to create the mage entity.
type MageComponents struct {
	Positions    *ecs.Storage[component.Position]
	Velocities   *ecs.Storage[component.Velocity]
	Facings      *ecs.Storage[component.Facing]
	Animations   *ecs.Storage[component.Animation]
	Directionals *ecs.Storage[component.DirectionalAnimation]
	SpriteSheets *ecs.Storage[component.SpriteSheet]
}

// MageFactoryConfig contains parameters for creating the mage.
type MageFactoryConfig struct {
	X, Y        float64
	SpriteSheet *ebiten.Image
	FrameWidth  int
	FrameHeight int
	Columns     int
	AnimFPS     float64
	SizeInUnits float64 // Target width in world units; 0 = native size
}

// CreateMage spawns a mage: a generic character (position, movement, facing, animation, sprite).
// Whether this mage is the one the player controls or the camera follows is the scene's call, so
// the factory leaves InputControlled and CameraTarget for the scene to attach.
func CreateMage(w *ecs.World, c *MageComponents, cfg MageFactoryConfig) ecs.Entity {
	e := w.Spawn()

	c.Positions.Set(e, component.Position{X: cfg.X, Y: cfg.Y})
	c.Velocities.Set(e, component.Velocity{X: 0, Y: 0})
	c.Facings.Set(e, component.Facing{Direction: component.FacingDown, Walking: false})

	// Four-directional sheet: one 6-frame walk row per direction, idle is each row's first frame.
	spec := component.DirectionalSheetSpec{
		Directions: component.DirectionsFour,
		Rows: map[component.FacingDirection]int{
			component.FacingDown:  0,
			component.FacingLeft:  6,
			component.FacingUp:    12,
			component.FacingRight: 18,
		},
		IdleFrames: 1,
		WalkFrames: 6,
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
