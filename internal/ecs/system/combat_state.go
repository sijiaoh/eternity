package system

import (
	"math"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

// CombatStateSystem decides whether the player's party is in combat: true when any living enemy
// stands within Radius of the player. That single, observable rule — an enemy has breached the
// guard radius around the player — is what flips allies from trailing to fighting. The flag is
// derived fresh each frame rather than stored, so there is no state to fall out of sync.
type CombatStateSystem struct {
	factions  *ecs.Storage[component.Faction]
	positions *ecs.Storage[component.Position]
	player    ecs.Entity
	radius    float64
}

func NewCombatStateSystem(
	factions *ecs.Storage[component.Faction],
	positions *ecs.Storage[component.Position],
	player ecs.Entity,
	radius float64,
) *CombatStateSystem {
	return &CombatStateSystem{
		factions:  factions,
		positions: positions,
		player:    player,
		radius:    radius,
	}
}

// InCombat reports whether an enemy (any faction opposing the player's) stands within the radius of
// the player. It is false when the player is gone or unplaced — a downed player has no party to
// command.
func (s *CombatStateSystem) InCombat(w *ecs.World) bool {
	if !w.Alive(s.player) {
		return false
	}
	playerFaction, ok := s.factions.Get(s.player)
	if !ok {
		return false
	}
	playerPos, ok := s.positions.Get(s.player)
	if !ok {
		return false
	}

	inCombat := false
	s.factions.Each(func(e ecs.Entity, faction *component.Faction) {
		if inCombat || *faction == playerFaction || !w.Alive(e) {
			return
		}
		pos, ok := s.positions.Get(e)
		if !ok {
			return
		}
		if math.Hypot(pos.X-playerPos.X, pos.Y-playerPos.Y) <= s.radius {
			inCombat = true
		}
	})
	return inCombat
}
