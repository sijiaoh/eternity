//go:build test

package system

import (
	"testing"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

// targetingFixture wires a world with the storages AITargetingSystem needs and a single goblin
// follower at the origin, so each spec only has to place enemies and run the system.
type targetingFixture struct {
	world     *ecs.World
	factions  *ecs.Storage[component.Faction]
	positions *ecs.Storage[component.Position]
	aiFollows *ecs.Storage[component.AIFollow]
	sys       *AITargetingSystem
	goblin    ecs.Entity
}

func newTargetingFixture() *targetingFixture {
	w := ecs.NewWorld(10)
	factions := ecs.NewStorage[component.Faction](10)
	positions := ecs.NewStorage[component.Position](10)
	aiFollows := ecs.NewStorage[component.AIFollow](10)

	goblin := w.Spawn()
	positions.Set(goblin, component.Position{X: 0, Y: 0})
	factions.Set(goblin, component.FactionEnemy)
	aiFollows.Set(goblin, component.AIFollow{Speed: 3})

	return &targetingFixture{
		world:     w,
		factions:  factions,
		positions: positions,
		aiFollows: aiFollows,
		sys:       NewAITargetingSystem(aiFollows, factions, positions),
		goblin:    goblin,
	}
}

func (f *targetingFixture) addEnemy(x, y float64) ecs.Entity {
	e := f.world.Spawn()
	f.positions.Set(e, component.Position{X: x, Y: y})
	f.factions.Set(e, component.FactionPlayer)
	return e
}

func (f *targetingFixture) target() ecs.Entity {
	ai, _ := f.aiFollows.Get(f.goblin)
	return ai.Target
}

func TestAITargeting_PicksNearestEnemy(t *testing.T) {
	f := newTargetingFixture()
	f.addEnemy(10, 0) // far
	near := f.addEnemy(3, 0)

	f.sys.Update(f.world, 0)

	if got := f.target(); got != near {
		t.Errorf("target = %v, want nearest enemy %v", got, near)
	}
}

func TestAITargeting_RetargetsWhenNearestDies(t *testing.T) {
	f := newTargetingFixture()
	far := f.addEnemy(10, 0)
	near := f.addEnemy(3, 0)

	f.sys.Update(f.world, 0)
	if got := f.target(); got != near {
		t.Fatalf("target = %v, want nearest enemy %v", got, near)
	}

	f.world.Despawn(near)
	f.sys.Update(f.world, 0)

	if got := f.target(); got != far {
		t.Errorf("after nearest died, target = %v, want next-nearest %v", got, far)
	}
}

func TestAITargeting_NoEnemyClearsTarget(t *testing.T) {
	f := newTargetingFixture()
	enemy := f.addEnemy(3, 0)

	f.sys.Update(f.world, 0)
	if got := f.target(); got != enemy {
		t.Fatalf("target = %v, want enemy %v", got, enemy)
	}

	f.world.Despawn(enemy)
	f.sys.Update(f.world, 0)

	if got := f.target(); f.world.Alive(got) {
		t.Errorf("with no enemies, target = %v (alive), want invalid/none", got)
	}
}

func TestAITargeting_IgnoresSameFaction(t *testing.T) {
	f := newTargetingFixture()

	// A closer ally of the goblin's own faction must not be chosen.
	friend := f.world.Spawn()
	f.positions.Set(friend, component.Position{X: 1, Y: 0})
	f.factions.Set(friend, component.FactionEnemy)

	enemy := f.addEnemy(5, 0)

	f.sys.Update(f.world, 0)

	if got := f.target(); got != enemy {
		t.Errorf("target = %v, want opposing-faction enemy %v (not same-faction friend)", got, enemy)
	}
}
