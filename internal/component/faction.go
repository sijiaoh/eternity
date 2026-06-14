package component

// Faction identifies which side a combat unit belongs to — the basis for
// friend/foe target selection.
type Faction int

const (
	// FactionPlayer covers the player-controlled Mage and its allies.
	FactionPlayer Faction = iota
	// FactionEnemy covers the goblins.
	FactionEnemy
)
