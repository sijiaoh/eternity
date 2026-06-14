package component

import "fmt"

// AnimationState defines a named animation sequence within a sprite sheet.
// StartFrame is the index of the first frame; subsequent frames are contiguous.
type AnimationState struct {
	Name       string
	StartFrame int     // first frame index in the sprite sheet
	FrameCount int     // number of frames in this animation
	FPS        float64 // playback speed in frames per second
	Loop       bool    // whether to loop after reaching the last frame
}

// Animation holds animation playback state for sprite sheet rendering.
type Animation struct {
	States       map[string]AnimationState
	CurrentState string
	FrameIndex   int     // current frame within the animation (0 to FrameCount-1)
	Elapsed      float64 // time elapsed since last frame change
	Finished     bool    // true when non-looping animation has ended
}

// NewAnimation creates an Animation with the given states.
// The first state in the slice becomes the initial state.
func NewAnimation(states []AnimationState) *Animation {
	a := &Animation{
		States: make(map[string]AnimationState),
	}
	for _, s := range states {
		a.States[s.Name] = s
	}
	if len(states) > 0 {
		a.CurrentState = states[0].Name
	}
	return a
}

// SetState switches the animation to a named state, resetting playback to its first
// frame. Switching to the current state is a no-op so playback continues uninterrupted.
// An unknown state is a programming error (a state name the entity's spec never declared)
// and panics rather than silently doing nothing.
func (a *Animation) SetState(name string) {
	if a.CurrentState == name {
		return
	}
	if _, ok := a.States[name]; !ok {
		panic(fmt.Sprintf("animation: unknown state %q", name))
	}
	a.CurrentState = name
	a.FrameIndex = 0
	a.Elapsed = 0
	a.Finished = false
}

// Frame returns the current frame index in the sprite sheet.
// This is StartFrame + FrameIndex within the current animation.
func (a *Animation) Frame() int {
	state, ok := a.States[a.CurrentState]
	if !ok {
		return 0
	}
	return state.StartFrame + a.FrameIndex
}
