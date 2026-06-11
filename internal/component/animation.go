package component

// AnimationState defines a named animation sequence within a sprite sheet.
// StartFrame is the index of the first frame; subsequent frames are contiguous.
type AnimationState struct {
	Name       string
	StartFrame int     // first frame index in the sprite sheet
	FrameCount int     // number of frames in this animation
	FPS        float64 // playback speed in frames per second
	Loop       bool    // whether to loop after reaching the last frame
}

// Animation manages sprite sheet animation playback.
//
// Use NewAnimation to create with a list of states. Call Update each frame
// with deltaTime to advance the animation. Use Frame to get the current
// sprite sheet index for rendering.
type Animation struct {
	states       map[string]AnimationState
	currentState string
	frameIndex   int     // current frame within the animation (0 to FrameCount-1)
	elapsed      float64 // time elapsed since last frame change
	finished     bool    // true when non-looping animation has ended
}

// NewAnimation creates an Animation with the given states.
// The first state in the slice becomes the initial state.
func NewAnimation(states []AnimationState) *Animation {
	a := &Animation{
		states: make(map[string]AnimationState),
	}
	for _, s := range states {
		a.states[s.Name] = s
	}
	if len(states) > 0 {
		a.currentState = states[0].Name
	}
	return a
}

// Update advances the animation by deltaTime seconds.
func (a *Animation) Update(deltaTime float64) {
	state, ok := a.states[a.currentState]
	if !ok || state.FrameCount <= 1 || state.FPS <= 0 || a.finished {
		return
	}
	if deltaTime <= 0 {
		return
	}

	frameDuration := 1.0 / state.FPS
	a.elapsed += deltaTime

	for a.elapsed >= frameDuration {
		a.elapsed -= frameDuration
		a.frameIndex++

		if a.frameIndex >= state.FrameCount {
			if state.Loop {
				a.frameIndex = 0
			} else {
				a.frameIndex = state.FrameCount - 1
				a.finished = true
				return
			}
		}
	}
}

// Frame returns the current frame index in the sprite sheet.
// This is StartFrame + frameIndex within the current animation.
func (a *Animation) Frame() int {
	state, ok := a.states[a.currentState]
	if !ok {
		return 0
	}
	return state.StartFrame + a.frameIndex
}

// SetState switches to a different animation state.
// If the state is the same as current, nothing happens.
// If the state doesn't exist, nothing happens.
func (a *Animation) SetState(name string) {
	if a.currentState == name {
		return
	}
	if _, ok := a.states[name]; !ok {
		return
	}
	a.currentState = name
	a.frameIndex = 0
	a.elapsed = 0
	a.finished = false
}

func (a *Animation) State() string {
	return a.currentState
}

func (a *Animation) IsFinished() bool {
	return a.finished
}

// Reset restarts the current animation from the beginning.
func (a *Animation) Reset() {
	a.frameIndex = 0
	a.elapsed = 0
	a.finished = false
}
