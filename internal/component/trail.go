package component

import "eternity/internal/ecs"

// Trail marks an ally that follows a leader in single file, settling a fixed Gap behind it rather
// than overlapping. Chaining each Trail's Leader to the unit ahead forms a Dragon-Quest-style line
// (the scene builds it; see linkTrailFormation). Unlike AIFollow (whose Target is the nearest
// enemy, re-picked every frame, and which closes to overlap), Trail's Leader is fixed by the scene
// and the follower halts once within Gap.
type Trail struct {
	Leader ecs.Entity // unit ahead in the line to trail behind
	Speed  float64    // movement speed in units per second
	Gap    float64    // distance to keep behind the leader, in units
}
