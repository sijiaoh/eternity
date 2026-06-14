package scene

import (
	"testing"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

// allySpeed must stay 1.2× the player's speed, derived rather than hard-coded, so tuning the
// player's speed carries the allies with it.
func TestAllySpeed_DerivedFromPlayer(t *testing.T) {
	if got, want := allySpeed(), playerSpeed*1.2; got != want {
		t.Errorf("allySpeed() = %v, want %v", got, want)
	}
}

// linkTrailFormation must chain followers single-file: the head trails the leader and each
// subsequent ally trails the one ahead, with the shared trailing speed and gap.
func TestLinkTrailFormation_ChainsBehindLeader(t *testing.T) {
	w := ecs.NewWorld(10)
	trails := ecs.NewStorage[component.Trail](10)

	leader := w.Spawn()
	first := w.Spawn()
	second := w.Spawn()

	linkTrailFormation(trails, leader, []ecs.Entity{first, second})

	firstTrail, ok := trails.Get(first)
	if !ok || firstTrail.Leader != leader {
		t.Errorf("first.Leader = %v, want leader %v", firstTrail.Leader, leader)
	}
	secondTrail, ok := trails.Get(second)
	if !ok || secondTrail.Leader != first {
		t.Errorf("second.Leader = %v, want first %v", secondTrail.Leader, first)
	}
	if firstTrail.Speed != allySpeed() || firstTrail.Gap != trailGap || firstTrail.SlowRadius != trailSlowRadius {
		t.Errorf("first trail tuning = (speed %v, gap %v, slowRadius %v), want (%v, %v, %v)",
			firstTrail.Speed, firstTrail.Gap, firstTrail.SlowRadius, allySpeed(), trailGap, trailSlowRadius)
	}

	// The leader itself is just the anchor — it gets no Trail of its own.
	if _, ok := trails.Get(leader); ok {
		t.Error("leader should not receive a Trail component")
	}
}
