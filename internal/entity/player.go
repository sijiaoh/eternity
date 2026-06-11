package entity

import (
	"image"

	"ebiten-agent-example/internal/component"
	"ebiten-agent-example/internal/config"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	playerSpeed = 5.0 // units per second
)

type Facing int

const (
	FacingDown Facing = iota
	FacingLeft
	FacingUp
	FacingRight
)

type Player struct {
	Position  *component.Position
	Movement  *component.Movement
	Sprite    *component.Sprite
	Animation *component.Animation

	spriteSheet *ebiten.Image
	frameWidth  int
	frameHeight int
	columns     int

	facing  Facing
	walking bool
}

type PlayerConfig struct {
	SpriteSheet *ebiten.Image
	FrameWidth  int
	FrameHeight int
	Columns     int // number of columns in sprite sheet
	AnimFPS     float64
}

func NewPlayer(x, y float64, cfg PlayerConfig) *Player {
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

	return &Player{
		Position:    component.NewPosition(x, y),
		Movement:    component.NewMovement(playerSpeed),
		Sprite:      component.NewSprite(nil),
		Animation:   component.NewAnimation(states),
		spriteSheet: cfg.SpriteSheet,
		frameWidth:  cfg.FrameWidth,
		frameHeight: cfg.FrameHeight,
		columns:     cfg.Columns,
		facing:      FacingDown,
		walking:     false,
	}
}

func (p *Player) Update(deltaTime float64) {
	dir := component.MoveDirection{
		Up:    ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp),
		Down:  ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown),
		Left:  ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft),
		Right: ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight),
	}

	p.updateFacing(dir)
	p.Movement.Update(p.Position, dir, deltaTime)
	p.updateAnimationState()
	p.Animation.Update(deltaTime)
}

func (p *Player) updateFacing(dir component.MoveDirection) {
	p.walking = dir.Up || dir.Down || dir.Left || dir.Right

	// Priority: horizontal over vertical (common in 2D games)
	if dir.Left {
		p.facing = FacingLeft
	} else if dir.Right {
		p.facing = FacingRight
	} else if dir.Up {
		p.facing = FacingUp
	} else if dir.Down {
		p.facing = FacingDown
	}
	// If no direction pressed, keep the current facing
}

func (p *Player) updateAnimationState() {
	var stateName string
	prefix := "idle_"
	if p.walking {
		prefix = "walk_"
	}

	switch p.facing {
	case FacingDown:
		stateName = prefix + "down"
	case FacingLeft:
		stateName = prefix + "left"
	case FacingRight:
		stateName = prefix + "right"
	case FacingUp:
		stateName = prefix + "up"
	}

	p.Animation.SetState(stateName)
}

func (p *Player) Draw(screen *ebiten.Image) {
	frame := p.getFrame(p.Animation.Frame())
	pixelX := config.UnitsToPixels(p.Position.X)
	pixelY := config.UnitsToPixels(p.Position.Y)
	p.Sprite.DrawFrame(screen, frame, pixelX, pixelY)
}

func (p *Player) getFrame(index int) *ebiten.Image {
	if p.spriteSheet == nil {
		return nil
	}

	col := index % p.columns
	row := index / p.columns

	x := col * p.frameWidth
	y := row * p.frameHeight

	rect := image.Rect(x, y, x+p.frameWidth, y+p.frameHeight)
	return p.spriteSheet.SubImage(rect).(*ebiten.Image)
}
