package component

type MoveDirection struct {
	Up, Down, Left, Right bool
}

type Movement struct {
	Speed float64 // units per second
}

func NewMovement(speed float64) *Movement {
	return &Movement{Speed: speed}
}

// Update applies direction input to the position.
func (m *Movement) Update(pos *Position, dir MoveDirection, deltaTime float64) {
	dx, dy := m.CalcDelta(dir, deltaTime)
	pos.Move(dx, dy)
}

// CalcDelta returns movement delta for this frame.
// Diagonal movement is normalized to prevent faster diagonal speed.
func (m *Movement) CalcDelta(dir MoveDirection, deltaTime float64) (dx, dy float64) {
	if dir.Up {
		dy -= 1
	}
	if dir.Down {
		dy += 1
	}
	if dir.Left {
		dx -= 1
	}
	if dir.Right {
		dx += 1
	}

	if dx != 0 && dy != 0 {
		const diagonalFactor = 0.7071067811865476 // 1/sqrt(2)
		dx *= diagonalFactor
		dy *= diagonalFactor
	}

	dx *= m.Speed * deltaTime
	dy *= m.Speed * deltaTime

	return dx, dy
}
