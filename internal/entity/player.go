package entity

import (
	"image/color"

	"ebiten-agent-example/internal/component"
	"ebiten-agent-example/internal/input"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	playerRadius = 20
	playerSpeed  = 240.0 // pixels per second (≈4 pixels at 60fps)
)

type Player struct {
	Position component.Position
	Movement component.Movement
}

func NewPlayer(x, y float64) *Player {
	return &Player{
		Position: component.Position{X: x, Y: y},
		Movement: component.NewMovement(playerSpeed),
	}
}

func (p *Player) Update(deltaTime float64) {
	dir := input.GetDirection()
	moveDir := component.MoveDirection{
		Up:    dir.Up,
		Down:  dir.Down,
		Left:  dir.Left,
		Right: dir.Right,
	}
	dx, dy := p.Movement.CalcDelta(moveDir, deltaTime)
	p.Position.Move(dx, dy)
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
