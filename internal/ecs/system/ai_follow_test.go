//go:build test

package system

import (
	"math"
	"testing"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

func TestAIFollowSystem_MovesTowardTarget(t *testing.T) {
	w := ecs.NewWorld(10)
	positions := ecs.NewStorage[component.Position](10)
	velocities := ecs.NewStorage[component.Velocity](10)
	aiFollows := ecs.NewStorage[component.AIFollow](10)

	sys := NewAIFollowSystem(aiFollows, positions, velocities)

	target := w.Spawn()
	positions.Set(target, component.Position{X: 10, Y: 0})

	follower := w.Spawn()
	positions.Set(follower, component.Position{X: 0, Y: 0})
	velocities.Set(follower, component.Velocity{X: 0, Y: 0})
	aiFollows.Set(follower, component.AIFollow{Target: target, Speed: 5})

	sys.Update(w, 0)

	vel, _ := velocities.Get(follower)
	if vel.X != 5 || vel.Y != 0 {
		t.Errorf("velocity = (%v, %v), want (5, 0)", vel.X, vel.Y)
	}
}

func TestAIFollowSystem_NormalizesDirection(t *testing.T) {
	w := ecs.NewWorld(10)
	positions := ecs.NewStorage[component.Position](10)
	velocities := ecs.NewStorage[component.Velocity](10)
	aiFollows := ecs.NewStorage[component.AIFollow](10)

	sys := NewAIFollowSystem(aiFollows, positions, velocities)

	target := w.Spawn()
	positions.Set(target, component.Position{X: 3, Y: 4})

	follower := w.Spawn()
	positions.Set(follower, component.Position{X: 0, Y: 0})
	velocities.Set(follower, component.Velocity{X: 0, Y: 0})
	aiFollows.Set(follower, component.AIFollow{Target: target, Speed: 10})

	sys.Update(w, 0)

	vel, _ := velocities.Get(follower)
	// Direction is (3, 4), distance is 5, normalized is (0.6, 0.8), scaled by 10 is (6, 8)
	if math.Abs(vel.X-6) > 0.001 || math.Abs(vel.Y-8) > 0.001 {
		t.Errorf("velocity = (%v, %v), want (6, 8)", vel.X, vel.Y)
	}
}

func TestAIFollowSystem_StopsWhenAtTarget(t *testing.T) {
	w := ecs.NewWorld(10)
	positions := ecs.NewStorage[component.Position](10)
	velocities := ecs.NewStorage[component.Velocity](10)
	aiFollows := ecs.NewStorage[component.AIFollow](10)

	sys := NewAIFollowSystem(aiFollows, positions, velocities)

	target := w.Spawn()
	positions.Set(target, component.Position{X: 5, Y: 5})

	follower := w.Spawn()
	positions.Set(follower, component.Position{X: 5, Y: 5})
	velocities.Set(follower, component.Velocity{X: 10, Y: 10})
	aiFollows.Set(follower, component.AIFollow{Target: target, Speed: 5})

	sys.Update(w, 0)

	vel, _ := velocities.Get(follower)
	if vel.X != 0 || vel.Y != 0 {
		t.Errorf("velocity at target = (%v, %v), want (0, 0)", vel.X, vel.Y)
	}
}

func TestAIFollowSystem_StopsWhenTargetDead(t *testing.T) {
	w := ecs.NewWorld(10)
	positions := ecs.NewStorage[component.Position](10)
	velocities := ecs.NewStorage[component.Velocity](10)
	aiFollows := ecs.NewStorage[component.AIFollow](10)

	sys := NewAIFollowSystem(aiFollows, positions, velocities)

	target := w.Spawn()
	positions.Set(target, component.Position{X: 10, Y: 0})

	follower := w.Spawn()
	positions.Set(follower, component.Position{X: 0, Y: 0})
	velocities.Set(follower, component.Velocity{X: 5, Y: 0})
	aiFollows.Set(follower, component.AIFollow{Target: target, Speed: 5})

	w.Despawn(target)
	sys.Update(w, 0)

	vel, _ := velocities.Get(follower)
	if vel.X != 0 || vel.Y != 0 {
		t.Errorf("velocity with dead target = (%v, %v), want (0, 0)", vel.X, vel.Y)
	}
}

func TestAIFollowSystem_SkipsDeadFollower(t *testing.T) {
	w := ecs.NewWorld(10)
	positions := ecs.NewStorage[component.Position](10)
	velocities := ecs.NewStorage[component.Velocity](10)
	aiFollows := ecs.NewStorage[component.AIFollow](10)

	sys := NewAIFollowSystem(aiFollows, positions, velocities)

	target := w.Spawn()
	positions.Set(target, component.Position{X: 10, Y: 0})

	follower := w.Spawn()
	positions.Set(follower, component.Position{X: 0, Y: 0})
	velocities.Set(follower, component.Velocity{X: 0, Y: 0})
	aiFollows.Set(follower, component.AIFollow{Target: target, Speed: 5})

	w.Despawn(follower)
	sys.Update(w, 0)

	vel, _ := velocities.Get(follower)
	if vel.X != 0 || vel.Y != 0 {
		t.Errorf("dead follower velocity = (%v, %v), want (0, 0)", vel.X, vel.Y)
	}
}
