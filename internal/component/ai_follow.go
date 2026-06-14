package component

import "eternity/internal/ecs"

// AIFollow marks an entity that chases a target. AITargetingSystem sets Target each frame and
// AIFollowSystem moves toward it.
type AIFollow struct {
	Target ecs.Entity // entity to follow
	Speed  float64    // movement speed in units per second
}
