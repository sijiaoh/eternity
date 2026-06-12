package component

type FacingDirection int

const (
	FacingDown FacingDirection = iota
	FacingLeft
	FacingUp
	FacingRight
)

// Facing tracks entity facing direction and walking state.
type Facing struct {
	Direction FacingDirection
	Walking   bool
}

func NewFacing(direction FacingDirection) *Facing {
	return &Facing{Direction: direction, Walking: false}
}
