package scene

import (
	"testing"
	"testing/fstest"

	"eternity/internal/component"
	"eternity/internal/config"
	"eternity/internal/ecs"
	"eternity/internal/i18n"
	"eternity/internal/scenario"
)

func ptrF(v float64) *float64 { return &v }
func ptrB(v bool) *bool       { return &v }

// These specs pin how an optional scenario.Battle becomes the concrete situation the battle
// scene boots with: each field overrides when set, and falls back to the normal-play default
// when unset. They run headless, standing in for the ebiten-gated NewBattleScene.

func TestResolveMageStart(t *testing.T) {
	wantX := config.PixelsToUnits(config.ScreenWidth / 2)
	wantY := config.PixelsToUnits(config.ScreenHeight / 2)

	t.Run("defaults to screen center", func(t *testing.T) {
		x, y := resolveMageStart(scenario.Battle{})
		if x != wantX || y != wantY {
			t.Fatalf("got (%v, %v), want screen center (%v, %v)", x, y, wantX, wantY)
		}
	})

	t.Run("uses the scenario position when set", func(t *testing.T) {
		x, y := resolveMageStart(scenario.Battle{MageX: ptrF(1.5), MageY: ptrF(2.5)})
		if x != 1.5 || y != 2.5 {
			t.Fatalf("got (%v, %v), want (1.5, 2.5)", x, y)
		}
	})
}

func TestResolveGoblinStart(t *testing.T) {
	t.Run("defaults to an offset from the mage", func(t *testing.T) {
		x, y := resolveGoblinStart(scenario.Battle{}, 10, 20)
		if x != 10+goblinSpawnOffset || y != 20+goblinSpawnOffset {
			t.Fatalf("got (%v, %v), want mage+offset", x, y)
		}
	})

	t.Run("uses the scenario position when set", func(t *testing.T) {
		x, y := resolveGoblinStart(scenario.Battle{GoblinX: ptrF(7), GoblinY: ptrF(8)}, 10, 20)
		if x != 7 || y != 8 {
			t.Fatalf("got (%v, %v), want (7, 8)", x, y)
		}
	})
}

func TestResolveAllyStart(t *testing.T) {
	// The ally stands at a fixed offset to the player mage's left, at the same height. With no
	// movement driver attached, this spot is also where it stays — its standstill position.
	x, y := resolveAllyStart(10, 20)
	if x != 10-allySpawnOffset || y != 20 {
		t.Fatalf("got (%v, %v), want (%v, 20) — mage's left at same height", x, y, 10-allySpawnOffset)
	}
}

func TestAssignBattleFactions(t *testing.T) {
	world := ecs.NewWorld(8)
	factions := ecs.NewStorage[component.Faction](8)

	mage := world.Spawn()
	ally := world.Spawn()
	goblin := world.Spawn()

	assignBattleFactions(factions, []ecs.Entity{mage, ally}, []ecs.Entity{goblin})

	for _, tc := range []struct {
		name string
		e    ecs.Entity
		want component.Faction
	}{
		{"player mage", mage, component.FactionPlayer},
		{"ally mage", ally, component.FactionPlayer},
		{"goblin", goblin, component.FactionEnemy},
	} {
		got, ok := factions.Get(tc.e)
		if !ok {
			t.Errorf("%s: no faction assigned", tc.name)
			continue
		}
		if got != tc.want {
			t.Errorf("%s: faction = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestSpawnGoblin(t *testing.T) {
	if !spawnGoblin(scenario.Battle{}) {
		t.Error("unset goblin should default to spawning")
	}
	if spawnGoblin(scenario.Battle{Goblin: ptrB(false)}) {
		t.Error("goblin:false should suppress the goblin")
	}
	if !spawnGoblin(scenario.Battle{Goblin: ptrB(true)}) {
		t.Error("goblin:true should spawn the goblin")
	}
}

func TestResolveTimeScale(t *testing.T) {
	if got := resolveTimeScale(scenario.Battle{}); got != 1.0 {
		t.Errorf("unset = %v, want 1.0 (normal speed)", got)
	}
	if got := resolveTimeScale(scenario.Battle{TimeScale: ptrF(0)}); got != 0 {
		t.Errorf("0 = %v, want 0 (paused)", got)
	}
	if got := resolveTimeScale(scenario.Battle{TimeScale: ptrF(0.5)}); got != 0.5 {
		t.Errorf("0.5 = %v, want 0.5 (slow-mo)", got)
	}
}

func TestInitialDialogue(t *testing.T) {
	// A minimal bundle: sampleDialogue's keys may be missing, but Get echoes the key, so the
	// lines are still non-empty — enough to drive the active/inactive distinction we care about.
	bundle, err := i18n.Load(fstest.MapFS{"zh.json": {Data: []byte(`{}`)}}, "zh")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	t.Run("dialogue:true opens the scene already in dialogue", func(t *testing.T) {
		lines := initialDialogue(scenario.Battle{Dialogue: true}, bundle)
		if len(lines) == 0 {
			t.Fatal("expected initial dialogue lines, got none")
		}
		d := &component.Dialogue{}
		d.Start(lines)
		if !d.Active {
			t.Fatal("dialogue should be active when the scene starts in dialogue")
		}
	})

	t.Run("default leaves the scene without dialogue", func(t *testing.T) {
		lines := initialDialogue(scenario.Battle{}, bundle)
		if lines != nil {
			t.Fatalf("expected no initial dialogue, got %v", lines)
		}
		d := &component.Dialogue{}
		d.Start(lines)
		if d.Active {
			t.Fatal("dialogue should be inactive by default")
		}
	})
}
