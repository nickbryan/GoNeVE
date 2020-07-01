package engine

import (
	"log"
	"math"
	"sync"
	"time"
)

const defaultSimulationsPerSecond = 20

// SimulationStepper is responsible for handling a single step within the game simulation.
type SimulationStepper interface {
	Step(s Simulator, r Renderer, elapsed time.Duration)
}

// fixedStepSimulation allows world rendering to happen as fast as possible (when vSync is disabled) but limits the
// number of simulations to the specified limit (sps).
type fixedStepSimulation struct {
	initializer     sync.Once
	deltaTime       time.Duration
	frameStart      time.Duration
	accumulator     time.Duration
	secondsStart    time.Duration
	sps             uint
	maxRenderSkips  uint
	updates, frames int
}

// Step moves the simulation forward by one frame. If the simulation is running behind then multiple simulation frames
// may occur before rendering and ending the step. Render will be passed the current distance in time between the last update
// and the next update, as a decimal, so that interpolation can be calculated.
func (fss *fixedStepSimulation) Step(s Simulator, r Renderer, elapsed time.Duration) {
	fss.initializer.Do(func() {
		fss.secondsStart = elapsed
		fss.deltaTime = time.Second / time.Duration(fss.sps)
		fss.frameStart = elapsed
		fss.maxRenderSkips = uint(math.Ceil(float64(fss.sps) * 0.2)) // TODO: is this overkill?
	})

	frameTime := elapsed - fss.frameStart
	fss.accumulator += frameTime

	var loops uint = 0
	for fss.accumulator >= fss.deltaTime && loops < fss.maxRenderSkips {
		s.Simulate(fss.deltaTime.Seconds())

		fss.accumulator -= fss.deltaTime
		loops++
		fss.updates++
	}

	alpha := float64(fss.accumulator) / float64(fss.deltaTime)
	r.Render(alpha)

	if elapsed-fss.secondsStart >= time.Second {
		log.Printf("Fps: %d UPS: %d", fss.frames, fss.updates)
		fss.updates = 0
		fss.frames = 0
		fss.secondsStart = elapsed
	}

	fss.frames++
	fss.frameStart = elapsed
}
