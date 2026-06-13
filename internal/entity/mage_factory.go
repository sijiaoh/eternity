//go:build !test

package entity

import (
	"eternity/internal/component"
	"eternity/internal/ecs"

	"github.com/hajimehoshi/ebiten/v2"
)

// MageComponents holds all storages needed to create the mage entity.
type MageComponents struct {
	Positions     *ecs.Storage[component.Position]
	Velocities    *ecs.Storage[component.Velocity]
	Inputs        *ecs.Storage[component.InputControlled]
	Facings       *ecs.Storage[component.Facing]
	Animations    *ecs.Storage[component.Animation]
	SpriteSheets  *ecs.Storage[component.SpriteSheet]
	CameraTargets *ecs.Storage[component.CameraTarget]
}

// MageFactoryConfig contains parameters for creating the mage.
type MageFactoryConfig struct {
	X, Y        float64
	Speed       float64
	SpriteSheet *ebiten.Image
	FrameWidth  int
	FrameHeight int
	Columns     int
	AnimFPS     float64
	SizeInUnits float64 // Target width in world units; 0 = native size
}

// CreateMage spawns the mage entity (the player-controlled character) with all required components.
func CreateMage(w *ecs.World, c *MageComponents, cfg MageFactoryConfig) ecs.Entity {
	e := w.Spawn()

	c.Positions.Set(e, component.Position{X: cfg.X, Y: cfg.Y})
	c.Velocities.Set(e, component.Velocity{X: 0, Y: 0})
	c.Inputs.Set(e, component.InputControlled{Speed: cfg.Speed})
	c.Facings.Set(e, component.Facing{Direction: component.FacingDown, Walking: false})

	states := []component.AnimationState{
		{Name: "idle_down", StartFrame: 0, FrameCount: 1, FPS: cfg.AnimFPS, Loop: true},
		{Name: "walk_down", StartFrame: 0, FrameCount: 6, FPS: cfg.AnimFPS, Loop: true},
		{Name: "idle_left", StartFrame: 6, FrameCount: 1, FPS: cfg.AnimFPS, Loop: true},
		{Name: "walk_left", StartFrame: 6, FrameCount: 6, FPS: cfg.AnimFPS, Loop: true},
		{Name: "idle_up", StartFrame: 12, FrameCount: 1, FPS: cfg.AnimFPS, Loop: true},
		{Name: "walk_up", StartFrame: 12, FrameCount: 6, FPS: cfg.AnimFPS, Loop: true},
		{Name: "idle_right", StartFrame: 18, FrameCount: 1, FPS: cfg.AnimFPS, Loop: true},
		{Name: "walk_right", StartFrame: 18, FrameCount: 6, FPS: cfg.AnimFPS, Loop: true},
	}
	c.Animations.Set(e, *component.NewAnimation(states))

	c.SpriteSheets.Set(e, component.SpriteSheet{
		Image:       cfg.SpriteSheet,
		FrameWidth:  cfg.FrameWidth,
		FrameHeight: cfg.FrameHeight,
		Columns:     cfg.Columns,
		Anchor:      component.AnchorCenter(),
		SizeInUnits: cfg.SizeInUnits,
	})

	c.CameraTargets.Set(e, component.CameraTarget{})

	return e
}
