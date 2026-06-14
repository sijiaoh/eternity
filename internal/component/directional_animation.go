package component

import "fmt"

// AnimationStateName is the single source of truth for state names
// ("<motion>_<direction>", e.g. "walk_left"): DirectionalSheetSpec generates states from
// it and ResolveState looks them up by it, so generation and resolution can never disagree.
func AnimationStateName(walking bool, dir FacingDirection) string {
	motion := "idle"
	if walking {
		motion = "walk"
	}
	return motion + "_" + directionName(dir)
}

func directionName(dir FacingDirection) string {
	switch dir {
	case FacingDown:
		return "down"
	case FacingLeft:
		return "left"
	case FacingUp:
		return "up"
	case FacingRight:
		return "right"
	default:
		panic(fmt.Sprintf("directional animation: unknown facing direction %d", dir))
	}
}

// DirectionSet is the set of directions a character actually has artwork for.
// A character declares one; callers never need to know which — ResolveState maps
// any facing onto a state the set provides.
type DirectionSet struct {
	directions []FacingDirection
}

var (
	// DirectionsFour: distinct artwork for up/down/left/right.
	DirectionsFour = DirectionSet{[]FacingDirection{FacingDown, FacingLeft, FacingUp, FacingRight}}
	// DirectionsHorizontal: only left/right artwork; vertical movement reuses the
	// last horizontal facing (see ResolveState).
	DirectionsHorizontal = DirectionSet{[]FacingDirection{FacingLeft, FacingRight}}
)

// All returns the set's directions in a stable order (used to build state lists deterministically).
func (s DirectionSet) All() []FacingDirection { return s.directions }

func (s DirectionSet) has(dir FacingDirection) bool {
	for _, d := range s.directions {
		if d == dir {
			return true
		}
	}
	return false
}

// DirectionalAnimation is the direction→state resolution layer that sits between the
// generic Animation (pure frame playback) and the entity's Facing. It owns which
// directions the character supports and remembers the last horizontal facing so a
// horizontal-only character can keep facing sideways while moving up or down.
type DirectionalAnimation struct {
	Directions     DirectionSet
	lastHorizontal FacingDirection // FacingLeft or FacingRight
}

// NewDirectionalAnimation creates the resolution layer for a character with the given
// direction set. The initial horizontal facing is right; it is overwritten the first
// time the character faces left or right.
func NewDirectionalAnimation(dirs DirectionSet) *DirectionalAnimation {
	return &DirectionalAnimation{Directions: dirs, lastHorizontal: FacingRight}
}

// ResolveState returns the animation state name to play for the given motion and
// facing, honouring the character's direction set. For DirectionsFour the facing is
// used as-is. For DirectionsHorizontal, vertical facings fall back to the last
// horizontal facing, so the resolved direction is always one the set provides.
func (d *DirectionalAnimation) ResolveState(walking bool, facing FacingDirection) string {
	if facing == FacingLeft || facing == FacingRight {
		d.lastHorizontal = facing
	}
	dir := facing
	if !d.Directions.has(dir) {
		dir = d.lastHorizontal
	}
	return AnimationStateName(walking, dir)
}

// DirectionalSheetSpec declares, in one place, how a character's directional states
// map onto its sprite sheet. Generating the state list and the direction set from a
// single spec keeps them from ever being declared inconsistently.
type DirectionalSheetSpec struct {
	Directions DirectionSet            // directions this character supports
	Rows       map[FacingDirection]int // direction → first frame of that direction's row
	IdleFrames int                     // frame count for idle states
	WalkFrames int                     // frame count for walk states
	FPS        float64
}

// States builds the idle/walk AnimationState list for every direction in the set.
// Idle and walk share a row's start frame (idle is the row's first frame). A direction
// in the set without a Rows entry is a programming error and panics.
func (spec DirectionalSheetSpec) States() []AnimationState {
	states := make([]AnimationState, 0, len(spec.Directions.All())*2)
	for _, dir := range spec.Directions.All() {
		start, ok := spec.Rows[dir]
		if !ok {
			panic(fmt.Sprintf("DirectionalSheetSpec: no row for direction %d", dir))
		}
		states = append(states,
			AnimationState{Name: AnimationStateName(false, dir), StartFrame: start, FrameCount: spec.IdleFrames, FPS: spec.FPS, Loop: true},
			AnimationState{Name: AnimationStateName(true, dir), StartFrame: start, FrameCount: spec.WalkFrames, FPS: spec.FPS, Loop: true},
		)
	}
	return states
}
