package component

import "eternity/internal/ecs"

// AIFollow marks an entity to follow a target entity directly.
// Entities with this component will have their velocity updated by AIFollowSystem.
type AIFollow struct {
	Target ecs.Entity // entity to follow
	Speed  float64    // movement speed in units per second
}

func NewAIFollow(target ecs.Entity, speed float64) *AIFollow {
	return &AIFollow{Target: target, Speed: speed}
}
