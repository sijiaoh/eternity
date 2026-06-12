package component

type FacingDirection int

const (
	FacingDown FacingDirection = iota
	FacingLeft
	FacingUp
	FacingRight
)

// Facing tracks entity facing direction and walking state.
// Used by systems to determine animation state.
type Facing struct {
	Direction FacingDirection
	Walking   bool
}

func NewFacing(direction FacingDirection) *Facing {
	return &Facing{Direction: direction, Walking: false}
}

// AnimationStateName returns the animation state name based on current facing and walking.
func (f *Facing) AnimationStateName() string {
	prefix := "idle_"
	if f.Walking {
		prefix = "walk_"
	}

	switch f.Direction {
	case FacingDown:
		return prefix + "down"
	case FacingLeft:
		return prefix + "left"
	case FacingRight:
		return prefix + "right"
	case FacingUp:
		return prefix + "up"
	default:
		return prefix + "down"
	}
}
