//go:build test

package system

import (
	"testing"

	"eternity/internal/component"
	"eternity/internal/ecs"
)

type animStateFixture struct {
	world      *ecs.World
	sys        *AnimationStateSystem
	entity     ecs.Entity
	animations *ecs.Storage[component.Animation]
	facings    *ecs.Storage[component.Facing]
}

// newAnimStateFixture builds a world with one four-directional entity, wired the same way
// the factories do (Animation + DirectionalAnimation + Facing).
func newAnimStateFixture() animStateFixture {
	w := ecs.NewWorld(10)
	animations := ecs.NewStorage[component.Animation](10)
	facings := ecs.NewStorage[component.Facing](10)
	directionals := ecs.NewStorage[component.DirectionalAnimation](10)

	spec := component.DirectionalSheetSpec{
		Directions: component.DirectionsFour,
		Rows: map[component.FacingDirection]int{
			component.FacingDown: 0, component.FacingLeft: 6, component.FacingUp: 12, component.FacingRight: 18,
		},
		IdleFrames: 1, WalkFrames: 6, FPS: 8,
	}

	e := w.Spawn()
	animations.Set(e, *component.NewAnimation(spec.States()))
	directionals.Set(e, *component.NewDirectionalAnimation(spec.Directions))
	facings.Set(e, component.Facing{Direction: component.FacingDown, Walking: false})

	return animStateFixture{
		world:      w,
		sys:        NewAnimationStateSystem(animations, facings, directionals),
		entity:     e,
		animations: animations,
		facings:    facings,
	}
}

func (f animStateFixture) face(dir component.FacingDirection, walking bool) {
	f.facings.Set(f.entity, component.Facing{Direction: dir, Walking: walking})
}

// The system's own job is the wiring: read Facing, resolve through DirectionalAnimation,
// switch the Animation. Resolution rules and SetState semantics are unit-tested in the
// component package, so here we only check that facing flows through to the right state
// for both motion categories.
func TestAnimationStateSystem_FacingFlowsToState(t *testing.T) {
	f := newAnimStateFixture()

	f.face(component.FacingLeft, true)
	f.sys.Update(f.world, 0.016)
	if got := f.animations.GetPtr(f.entity).CurrentState; got != "walk_left" {
		t.Errorf("walking left: state = %s, want walk_left", got)
	}

	f.face(component.FacingDown, false)
	f.sys.Update(f.world, 0.016)
	if got := f.animations.GetPtr(f.entity).CurrentState; got != "idle_down" {
		t.Errorf("idle down: state = %s, want idle_down", got)
	}
}

func TestAnimationStateSystem_SkipsDeadEntities(t *testing.T) {
	f := newAnimStateFixture()
	f.face(component.FacingLeft, true)
	f.world.Despawn(f.entity)

	f.sys.Update(f.world, 0.016) // must not touch despawned entity

	if got := f.animations.GetPtr(f.entity).CurrentState; got == "walk_left" {
		t.Error("despawned entity should not have its animation state updated")
	}
}
