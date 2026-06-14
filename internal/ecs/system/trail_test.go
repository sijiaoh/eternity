//go:build test

package system

import (
	"math"
	"testing"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

// newTrailFixture wires the storages and system used across the trail tests.
func newTrailFixture() (*ecs.World, *ecs.Storage[component.Position], *ecs.Storage[component.Velocity], *ecs.Storage[component.Trail], *TrailSystem) {
	w := ecs.NewWorld(10)
	positions := ecs.NewStorage[component.Position](10)
	velocities := ecs.NewStorage[component.Velocity](10)
	trails := ecs.NewStorage[component.Trail](10)
	return w, positions, velocities, trails, NewTrailSystem(trails, positions, velocities)
}

// Beyond Gap+SlowRadius, the follower heads straight at its leader at exactly Speed: direction is
// normalized (so distance doesn't inflate the velocity) and then scaled by Speed.
func TestTrailSystem_MovesTowardLeaderAtSpeed(t *testing.T) {
	w, positions, velocities, trails, sys := newTrailFixture()

	leader := w.Spawn()
	positions.Set(leader, component.Position{X: 3, Y: 4}) // distance 5, beyond Gap+SlowRadius (3)

	follower := w.Spawn()
	positions.Set(follower, component.Position{X: 0, Y: 0})
	velocities.Set(follower, component.Velocity{})
	trails.Set(follower, component.Trail{Leader: leader, Speed: 10, Gap: 1, SlowRadius: 2})

	sys.Update(w, 0)

	// Direction (3,4) normalized is (0.6,0.8); scaled by speed 10 is (6,8).
	vel, _ := velocities.Get(follower)
	if math.Abs(vel.X-6) > 0.001 || math.Abs(vel.Y-8) > 0.001 {
		t.Errorf("velocity = (%v, %v), want (6, 8)", vel.X, vel.Y)
	}
}

// Within the SlowRadius band the desired speed ramps down linearly with distance, so the follower
// eases in rather than charging at full Speed and slamming to a halt. The sample sits a quarter of
// the way into the band (an asymmetric point) so the assertion pins the ramp as linear: an eased
// curve like smoothstep would also pass at the midpoint but not here.
func TestTrailSystem_DeceleratesWithinSlowRadius(t *testing.T) {
	w, positions, velocities, trails, sys := newTrailFixture()

	leader := w.Spawn()
	positions.Set(leader, component.Position{X: 1.5, Y: 0}) // 0.5 into the 2-unit band past the gap

	follower := w.Spawn()
	positions.Set(follower, component.Position{X: 0, Y: 0})
	velocities.Set(follower, component.Velocity{})
	trails.Set(follower, component.Trail{Leader: leader, Speed: 6, Gap: 1, SlowRadius: 2})

	sys.Update(w, 0)

	// speed = 6 * (1.5-1)/2 = 1.5, straight along +X. (smoothstep at t=0.25 would give ~0.94.)
	vel, _ := velocities.Get(follower)
	if math.Abs(vel.X-1.5) > 0.001 || math.Abs(vel.Y) > 0.001 {
		t.Errorf("velocity = (%v, %v), want (1.5, 0)", vel.X, vel.Y)
	}
}

// Just past the gap the ramped speed falls below trailStopSpeed, so the follower stops outright
// instead of creeping — otherwise the residual velocity would read as "walking" and the ally would
// march in place after the player halts.
func TestTrailSystem_StopsInsideDeadzone(t *testing.T) {
	w, positions, velocities, trails, sys := newTrailFixture()

	leader := w.Spawn()
	positions.Set(leader, component.Position{X: 1.01, Y: 0}) // barely past the gap

	follower := w.Spawn()
	positions.Set(follower, component.Position{X: 0, Y: 0})
	velocities.Set(follower, component.Velocity{X: 6, Y: 0})
	trails.Set(follower, component.Trail{Leader: leader, Speed: 6, Gap: 1, SlowRadius: 2})

	sys.Update(w, 0)

	// speed = 6 * 0.01/2 = 0.03 < trailStopSpeed, so velocity snaps to exactly zero.
	vel, _ := velocities.Get(follower)
	if vel.X != 0 || vel.Y != 0 {
		t.Errorf("velocity inside deadzone = (%v, %v), want (0, 0)", vel.X, vel.Y)
	}
}

func TestTrailSystem_HaltsWithinGap(t *testing.T) {
	w, positions, velocities, trails, sys := newTrailFixture()

	leader := w.Spawn()
	positions.Set(leader, component.Position{X: 1, Y: 0}) // 1 unit away, exactly the gap

	follower := w.Spawn()
	positions.Set(follower, component.Position{X: 0, Y: 0})
	velocities.Set(follower, component.Velocity{X: 6, Y: 0})
	trails.Set(follower, component.Trail{Leader: leader, Speed: 6, Gap: 1})

	sys.Update(w, 0)

	// At/within the gap the follower keeps its distance and stops, so it stays behind the leader.
	vel, _ := velocities.Get(follower)
	if vel.X != 0 || vel.Y != 0 {
		t.Errorf("velocity within gap = (%v, %v), want (0, 0)", vel.X, vel.Y)
	}
}

func TestTrailSystem_StopsWhenLeaderDead(t *testing.T) {
	w, positions, velocities, trails, sys := newTrailFixture()

	leader := w.Spawn()
	positions.Set(leader, component.Position{X: 10, Y: 0})

	follower := w.Spawn()
	positions.Set(follower, component.Position{X: 0, Y: 0})
	velocities.Set(follower, component.Velocity{X: 6, Y: 0})
	trails.Set(follower, component.Trail{Leader: leader, Speed: 6, Gap: 1})

	w.Despawn(leader)
	sys.Update(w, 0)

	vel, _ := velocities.Get(follower)
	if vel.X != 0 || vel.Y != 0 {
		t.Errorf("velocity with dead leader = (%v, %v), want (0, 0)", vel.X, vel.Y)
	}
}

// TestTrailSystem_ChainFollowsImmediateLeader is the Dragon-Quest line spec: each follower trails
// the unit directly ahead of it, not the head. The three units sit at a right angle so "follow
// first" and "follow leader" point in different directions, letting the assertions tell them
// apart — if `second` wrongly chased the leader, its velocity would gain a Y component.
func TestTrailSystem_ChainFollowsImmediateLeader(t *testing.T) {
	w, positions, velocities, trails, sys := newTrailFixture()

	leader := w.Spawn()
	positions.Set(leader, component.Position{X: 0, Y: 0})

	first := w.Spawn()
	positions.Set(first, component.Position{X: 0, Y: 3}) // straight below the leader
	velocities.Set(first, component.Velocity{})
	trails.Set(first, component.Trail{Leader: leader, Speed: 6, Gap: 1})

	second := w.Spawn()
	positions.Set(second, component.Position{X: 3, Y: 3}) // beside first, away from the leader
	velocities.Set(second, component.Velocity{})
	trails.Set(second, component.Trail{Leader: first, Speed: 6, Gap: 1})

	sys.Update(w, 0)

	// first heads at the leader: straight -Y.
	firstVel, _ := velocities.Get(first)
	if math.Abs(firstVel.X) > 0.001 || firstVel.Y >= 0 {
		t.Errorf("first velocity = (%v, %v), want straight -Y toward the leader", firstVel.X, firstVel.Y)
	}
	// second heads at first (not the leader): straight -X, with no Y pull toward the head.
	secondVel, _ := velocities.Get(second)
	if secondVel.X >= 0 || math.Abs(secondVel.Y) > 0.001 {
		t.Errorf("second velocity = (%v, %v), want straight -X toward first", secondVel.X, secondVel.Y)
	}
}
