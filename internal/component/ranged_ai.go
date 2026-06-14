package component

// RangedAI marks an ally that fights at range: in combat it flees its nearest enemy, backing away
// to keep distance. It governs only the flee motion; staying near the player is the Leash's job, so
// the two constraints compose — RangedAI pushes the ally outward, Leash caps how far it gets.
// Compare AIFollow, which closes on its target instead of backing off.
type RangedAI struct {
	Speed float64 // backing-away speed in units per second
}
