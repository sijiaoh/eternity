package component

import "testing"

func TestHealth_IsDead(t *testing.T) {
	tests := []struct {
		name   string
		health Health
		want   bool
	}{
		{"zero health", Health{Current: 0, Max: 100}, true},
		{"negative health", Health{Current: -10, Max: 100}, true},
		{"positive health", Health{Current: 50, Max: 100}, false},
		{"full health", Health{Current: 100, Max: 100}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.health.IsDead(); got != tt.want {
				t.Errorf("IsDead() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealth_TakeDamage(t *testing.T) {
	tests := []struct {
		name        string
		initial     Health
		damage      int
		wantCurrent int
	}{
		{"normal damage", Health{Current: 100, Max: 100}, 30, 70},
		{"overkill damage", Health{Current: 50, Max: 100}, 100, 0},
		{"zero damage", Health{Current: 50, Max: 100}, 0, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.initial
			h.TakeDamage(tt.damage)
			if h.Current != tt.wantCurrent {
				t.Errorf("TakeDamage(%d) resulted in Current = %d, want %d", tt.damage, h.Current, tt.wantCurrent)
			}
		})
	}
}

func TestHealth_Heal(t *testing.T) {
	tests := []struct {
		name        string
		initial     Health
		amount      int
		wantCurrent int
	}{
		{"normal heal", Health{Current: 50, Max: 100}, 30, 80},
		{"overheal capped", Health{Current: 80, Max: 100}, 50, 100},
		{"zero heal", Health{Current: 50, Max: 100}, 0, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := tt.initial
			h.Heal(tt.amount)
			if h.Current != tt.wantCurrent {
				t.Errorf("Heal(%d) resulted in Current = %d, want %d", tt.amount, h.Current, tt.wantCurrent)
			}
		})
	}
}
