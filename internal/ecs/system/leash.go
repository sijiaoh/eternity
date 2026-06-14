package system

import (
	"math"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

// LeashSystem keeps each leashed unit within its anchor's radius. It runs after MovementSystem as a
// constraint resolver: a unit that integrated past the rim is snapped back onto it, and its
// outward velocity is cancelled so facing/animation read the along-rim slide — and the unit rests
// instead of walking in place when pinned dead against the rim. Inside the radius it does nothing,
// so the leash only bites during combat, when an AI would otherwise carry the ally off the player.
type LeashSystem struct {
	leashes    *ecs.Storage[component.Leash]
	positions  *ecs.Storage[component.Position]
	velocities *ecs.Storage[component.Velocity]
}

func NewLeashSystem(
	leashes *ecs.Storage[component.Leash],
	positions *ecs.Storage[component.Position],
	velocities *ecs.Storage[component.Velocity],
) *LeashSystem {
	return &LeashSystem{
		leashes:    leashes,
		positions:  positions,
		velocities: velocities,
	}
}

func (s *LeashSystem) Update(w *ecs.World, dt float64) {
	s.leashes.Each(func(e ecs.Entity, l *component.Leash) {
		if !w.Alive(e) {
			return
		}
		pos := s.positions.GetPtr(e)
		if pos == nil || !w.Alive(l.Anchor) {
			return
		}
		anchor, ok := s.positions.Get(l.Anchor)
		if !ok {
			return
		}

		dx := pos.X - anchor.X
		dy := pos.Y - anchor.Y
		dist := math.Hypot(dx, dy)
		if dist <= l.Range {
			return
		}

		// Snap back onto the rim along the anchor→unit direction.
		outX := dx / dist
		outY := dy / dist
		pos.X = anchor.X + outX*l.Range
		pos.Y = anchor.Y + outY*l.Range

		// Cancel the outward part of the velocity so it stops pushing past the rim; facing and
		// animation then reflect only the along-rim motion (and rest when the push is purely out).
		vel := s.velocities.GetPtr(e)
		if vel == nil {
			return
		}
		radial := vel.X*outX + vel.Y*outY
		if radial > 0 {
			vel.X -= radial * outX
			vel.Y -= radial * outY
		}
	})
}
