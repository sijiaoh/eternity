package component

type MoveDirection struct {
	Up, Down, Left, Right bool
}

type Movement struct {
	Speed float64
}

func NewMovement(speed float64) Movement {
	return Movement{Speed: speed}
}

// CalcDelta normalizes diagonal movement to prevent faster diagonal speed.
func (m *Movement) CalcDelta(dir MoveDirection) (dx, dy float64) {
	if dir.Up {
		dy -= m.Speed
	}
	if dir.Down {
		dy += m.Speed
	}
	if dir.Left {
		dx -= m.Speed
	}
	if dir.Right {
		dx += m.Speed
	}

	if dx != 0 && dy != 0 {
		const diagonalFactor = 0.7071067811865476 // 1/sqrt(2)
		dx *= diagonalFactor
		dy *= diagonalFactor
	}

	return dx, dy
}
