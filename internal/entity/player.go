package entity

import (
	"image/color"

	"ebiten-agent-example/internal/component"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	playerRadius = 20
	playerSpeed  = 240.0 // pixels per second (≈4 pixels at 60fps)
)

type Player struct {
	Position *component.Position
	Movement *component.Movement
}

func NewPlayer(x, y float64) *Player {
	return &Player{
		Position: component.NewPosition(x, y),
		Movement: component.NewMovement(playerSpeed),
	}
}

func (p *Player) Update(deltaTime float64) {
	dir := component.MoveDirection{
		Up:    ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp),
		Down:  ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown),
		Left:  ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft),
		Right: ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight),
	}
	p.Movement.Update(p.Position, dir, deltaTime)
}

func (p *Player) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(
		screen,
		float32(p.Position.X),
		float32(p.Position.Y),
		playerRadius,
		color.RGBA{R: 100, G: 149, B: 237, A: 255}, // cornflower blue
		true,
	)
}
