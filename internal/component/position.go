package component

type Position struct {
	X, Y float64
}

func (p *Position) Move(dx, dy float64) {
	p.X += dx
	p.Y += dy
}

type Velocity struct {
	X, Y float64
}
