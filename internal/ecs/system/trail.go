package system

import (
	"math"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

// trailStopSpeed is the speed, in units per second, below which a follower is treated as stopped:
// its velocity is zeroed instead of left as a creep. This keeps the arrival ramp's tail from
// leaving a tiny residual velocity that FacingSystem would read as "walking" — the cause of an ally
// marching in place after the player halts. It is the trailing counterpart to FacingSystem's
// "any non-zero velocity is walking" rule, so anything below it must reach exactly zero.
const trailStopSpeed = 0.1

// TrailSystem moves each Trail entity toward its leader using arrival steering: beyond the leader's
// Gap+SlowRadius it travels at full Speed, then eases down through the SlowRadius band to a stop at
// the Gap, so the follower settles behind it instead of overshooting and jittering. Chained leaders
// form a single-file line (see component.Trail). It reads the leader's current position, so run it
// before MovementSystem.
type TrailSystem struct {
	trails     *ecs.Storage[component.Trail]
	positions  *ecs.Storage[component.Position]
	velocities *ecs.Storage[component.Velocity]
}

func NewTrailSystem(
	trails *ecs.Storage[component.Trail],
	positions *ecs.Storage[component.Position],
	velocities *ecs.Storage[component.Velocity],
) *TrailSystem {
	return &TrailSystem{
		trails:     trails,
		positions:  positions,
		velocities: velocities,
	}
}

func (s *TrailSystem) Update(w *ecs.World, dt float64) {
	s.trails.Each(func(e ecs.Entity, t *component.Trail) {
		if !w.Alive(e) {
			return
		}

		pos := s.positions.GetPtr(e)
		vel := s.velocities.GetPtr(e)
		if pos == nil || vel == nil {
			return
		}

		if !w.Alive(t.Leader) {
			vel.X = 0
			vel.Y = 0
			return
		}

		leaderPos, ok := s.positions.Get(t.Leader)
		if !ok {
			vel.X = 0
			vel.Y = 0
			return
		}

		dx := leaderPos.X - pos.X
		dy := leaderPos.Y - pos.Y
		dist := math.Hypot(dx, dy)

		speed := arrivalSpeed(dist, t.Gap, t.SlowRadius, t.Speed)
		if speed == 0 {
			vel.X = 0
			vel.Y = 0
			return
		}

		vel.X = (dx / dist) * speed
		vel.Y = (dy / dist) * speed
	})
}

// arrivalSpeed returns the desired trailing speed at the given distance to the leader: 0 within Gap,
// a linear ramp from 0 to speed across the SlowRadius band, and full speed beyond it. Speeds below
// trailStopSpeed snap to 0 so the follower comes to a clean stop rather than creeping in place.
// dist > Gap whenever this returns non-zero, so callers can safely normalize by dist.
func arrivalSpeed(dist, gap, slowRadius, speed float64) float64 {
	if dist <= gap {
		return 0
	}
	desired := speed
	if slowRadius > 0 && dist < gap+slowRadius {
		desired = speed * (dist - gap) / slowRadius
	}
	if desired < trailStopSpeed {
		return 0
	}
	return desired
}
