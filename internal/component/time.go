package component

// Clock represents a time source that provides delta time for game logic.
// Clocks can be nested to create hierarchical time control.
//
// Root clocks track raw delta time directly. Child clocks derive their
// delta time from their parent, applying their own scale factor.
type Clock struct {
	scale     float64
	rawDelta  float64
	totalTime float64
	parent    *Clock
}

// NewClock creates a root clock with scale 1.0.
func NewClock() *Clock {
	return &Clock{scale: 1.0}
}

// NewChildClock creates a clock that derives time from a parent.
func NewChildClock(parent *Clock) *Clock {
	return &Clock{scale: 1.0, parent: parent}
}

// Update sets the raw delta time for this frame and advances TotalTime.
// For root clocks: pass the frame delta in seconds.
// For child clocks: pass any value (ignored); or use Tick() instead.
func (c *Clock) Update(rawDelta float64) {
	if c.parent == nil {
		c.rawDelta = rawDelta
	}
	c.totalTime += c.DeltaTime()
}

// Tick advances TotalTime without changing rawDelta. Use for child clocks.
func (c *Clock) Tick() {
	c.totalTime += c.DeltaTime()
}

// DeltaTime returns the scaled delta time for this frame.
// For child clocks, this recursively computes parent.DeltaTime() * scale.
func (c *Clock) DeltaTime() float64 {
	if c.parent != nil {
		return c.parent.DeltaTime() * c.scale
	}
	return c.rawDelta * c.scale
}

func (c *Clock) TotalTime() float64 {
	return c.totalTime
}

func (c *Clock) Scale() float64 {
	return c.scale
}

// SetScale sets the time scale. Negative values are clamped to 0.
func (c *Clock) SetScale(scale float64) {
	if scale < 0 {
		scale = 0
	}
	c.scale = scale
}

func (c *Clock) IsPaused() bool {
	return c.scale == 0
}

func (c *Clock) Pause() {
	c.scale = 0
}

func (c *Clock) Resume() {
	c.scale = 1.0
}
