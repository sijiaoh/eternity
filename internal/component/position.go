package component

type Position struct {
	X, Y float64
}

func NewPosition(x, y float64) *Position {
	return &Position{X: x, Y: y}
}

func (p *Position) Move(dx, dy float64) {
	p.X += dx
	p.Y += dy
}

type Velocity struct {
	X, Y float64
}

func NewVelocity(x, y float64) *Velocity {
	return &Velocity{X: x, Y: y}
}
