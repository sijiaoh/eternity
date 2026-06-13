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
	SpriteSheets *ecs.Storage[component.SpriteSheet]
}

// GoblinFactoryConfig contains parameters for creating a goblin.
type GoblinFactoryConfig struct {
	X, Y        float64
	Speed       float64
	Target      ecs.Entity
	SpriteSheet *ebiten.Image
	FrameWidth  int
	FrameHeight int
	Columns     int
	AnimFPS     float64
	SizeInUnits float64 // Target width in world units; 0 = native size
}

// CreateGoblin spawns a goblin entity with follow AI.
// Each 11-frame row is one facing direction (down/left/up/right), with frame 0
// as the idle pose — see sprite.source.md for the full sheet layout.
func CreateGoblin(w *ecs.World, c *GoblinComponents, cfg GoblinFactoryConfig) ecs.Entity {
	e := w.Spawn()

	c.Positions.Set(e, component.Position{X: cfg.X, Y: cfg.Y})
	c.Velocities.Set(e, component.Velocity{X: 0, Y: 0})
	c.AIFollows.Set(e, component.AIFollow{Target: cfg.Target, Speed: cfg.Speed})
	c.Facings.Set(e, component.Facing{Direction: component.FacingDown, Walking: false})

	states := []component.AnimationState{
		{Name: "idle_down", StartFrame: 0, FrameCount: 1, FPS: cfg.AnimFPS, Loop: true},
		{Name: "walk_down", StartFrame: 0, FrameCount: 8, FPS: cfg.AnimFPS, Loop: true},
		{Name: "idle_left", StartFrame: 11, FrameCount: 1, FPS: cfg.AnimFPS, Loop: true},
		{Name: "walk_left", StartFrame: 11, FrameCount: 8, FPS: cfg.AnimFPS, Loop: true},
		{Name: "idle_up", StartFrame: 22, FrameCount: 1, FPS: cfg.AnimFPS, Loop: true},
		{Name: "walk_up", StartFrame: 22, FrameCount: 8, FPS: cfg.AnimFPS, Loop: true},
		{Name: "idle_right", StartFrame: 33, FrameCount: 1, FPS: cfg.AnimFPS, Loop: true},
		{Name: "walk_right", StartFrame: 33, FrameCount: 8, FPS: cfg.AnimFPS, Loop: true},
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

	return e
}
