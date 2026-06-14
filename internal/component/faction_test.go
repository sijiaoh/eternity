package component

import "testing"

// Faction is a behaviorless marker enum; its only spec is that the two sides
// are distinct, which is what friend/foe targeting relies on.
func TestFactionSidesAreDistinct(t *testing.T) {
	if FactionPlayer == FactionEnemy {
		t.Fatal("player and enemy factions must be distinct values")
	}
}
