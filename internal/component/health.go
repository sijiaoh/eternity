package component

type Health struct {
	Current int
	Max     int
}

func (h *Health) IsDead() bool {
	return h.Current <= 0
}

func (h *Health) TakeDamage(damage int) {
	h.Current -= damage
	if h.Current < 0 {
		h.Current = 0
	}
}

func (h *Health) Heal(amount int) {
	h.Current += amount
	if h.Current > h.Max {
		h.Current = h.Max
	}
}
