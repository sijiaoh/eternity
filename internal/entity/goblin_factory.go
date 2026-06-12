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
}

// CreateGoblin spawns a goblin entity with follow AI.
// The goblin sprite sheet has only left-facing frames, so we use FlipOnRight
// to mirror the sprite when facing right. All directions map to the same frames.
func CreateGoblin(w *ecs.World, c *GoblinComponents, cfg GoblinFactoryConfig) ecs.Entity {
	e := w.Spawn()

	c.Positions.Set(e, component.Position{X: cfg.X, Y: cfg.Y})
	c.Velocities.Set(e, component.Velocity{X: 0, Y: 0})
	c.AIFollows.Set(e, component.AIFollow{Target: cfg.Target, Speed: cfg.Speed})
	c.Facings.Set(e, component.Facing{Direction: component.FacingLeft, Walking: false})

	// Goblin sprite sheet layout (10 columns, knife goblin = first 5 rows):
	// Row 0: Idle (4 frames)
	// Row 1: Gesture (4 frames)
	// Row 2: Walk (4 frames)
	// Row 3: Attack (6 frames)
	// Row 4: Death (4 frames)
	// All frames face left; FlipOnRight handles right direction.
	idleStart := 0 // Row 0
	walkStart := 20 // Row 2 (2 * 10 columns)

	states := []component.AnimationState{
		{Name: "idle_down", StartFrame: idleStart, FrameCount: 4, FPS: cfg.AnimFPS, Loop: true},
		{Name: "idle_left", StartFrame: idleStart, FrameCount: 4, FPS: cfg.AnimFPS, Loop: true},
		{Name: "idle_up", StartFrame: idleStart, FrameCount: 4, FPS: cfg.AnimFPS, Loop: true},
		{Name: "idle_right", StartFrame: idleStart, FrameCount: 4, FPS: cfg.AnimFPS, Loop: true},
		{Name: "walk_down", StartFrame: walkStart, FrameCount: 4, FPS: cfg.AnimFPS, Loop: true},
		{Name: "walk_left", StartFrame: walkStart, FrameCount: 4, FPS: cfg.AnimFPS, Loop: true},
		{Name: "walk_up", StartFrame: walkStart, FrameCount: 4, FPS: cfg.AnimFPS, Loop: true},
		{Name: "walk_right", StartFrame: walkStart, FrameCount: 4, FPS: cfg.AnimFPS, Loop: true},
	}
	c.Animations.Set(e, *component.NewAnimation(states))

	c.SpriteSheets.Set(e, component.SpriteSheet{
		Image:       cfg.SpriteSheet,
		FrameWidth:  cfg.FrameWidth,
		FrameHeight: cfg.FrameHeight,
		Columns:     cfg.Columns,
		Anchor:      component.AnchorCenter(),
		FlipOnRight: true,
	})

	return e
}
