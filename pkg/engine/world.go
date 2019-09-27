package engine

// Simulator represents a system that should be called as part of the main loop.
type Simulator interface {
	Simulate(dt float64)
}

// PreSimulator allows systems to be called before the main Simulator systems.
type PreSimulator interface {
	PreSimulate(dt float64)
}

// PreSimulator allows systems to be called after the main Simulator systems.
type PostSimulator interface {
	PostSimulate(dt float64)
}

// Renderer represents a system that has rendering functionality. In the default
// SimulationStepper this will be triggered separately to the different Simulators.
type Renderer interface {
	Render(ipl float64)
}

// PreRenderer allows render systems to be called before the main Renderer systems.
type PreRenderer interface {
	PreRender(ipl float64)
}

// PostRenderer allows render systems to be called after the main Renderer systems.
type PostRenderer interface {
	PostRender(ipl float64)
}

// World contains all of the different systems and is responsible for
// orchestrating the execution of the encapsulated systems.
type World struct {
	preSimulators  []PreSimulator
	postSimulators []PostSimulator
	simulators     []Simulator
	preRenderers   []PreRenderer
	postRenderers  []PostRenderer
	renderers      []Renderer
}

// NewWorld returns a new World instance.
func NewWorld() *World {
	return &World{}
}

// AddPreSimulator adds the given PreSimulator to the World.
func (w *World) AddPreSimulator(s PreSimulator) {
	w.preSimulators = append(w.preSimulators, s)
}

// AddPostSimulator adds the given PostSimulator to the World.
func (w *World) AddPostSimulator(s PostSimulator) {
	w.postSimulators = append(w.postSimulators, s)
}

// AddSimulator adds the given Simulator to the World.
func (w *World) AddSimulator(s Simulator) {
	w.simulators = append(w.simulators, s)
}

// AddRenderer adds the given Renderer to the World.
func (w *World) AddRenderer(r Renderer) {
	w.renderers = append(w.renderers, r)
}

// AddPreRenderer adds teh given PreRenderer to the World.
func (w *World) AddPreRenderer(r PreRenderer) {
	w.preRenderers = append(w.preRenderers, r)
}

// AddPostRenderer adds the given PostRenderer to the World.
func (w *World) AddPostRenderer(r PostRenderer) {
	w.postRenderers = append(w.postRenderers, r)
}

// Simulate will usually be called by the main game loop. It will loop through PreSimulators then
// the Simulators and finally the PreSimulators in a single call. All Simulators will be called in the
// order that they were added.
func (w *World) Simulate(dt float64) {
	for _, s := range w.preSimulators {
		s.PreSimulate(dt)
	}

	for _, s := range w.simulators {
		s.Simulate(dt)
	}

	for _, s := range w.postSimulators {
		s.PostSimulate(dt)
	}
}

// Render calls all of the Renderers in the order that they were added.
func (w *World) Render(ipl float64) {
	for _, r := range w.preRenderers {
		r.PreRender(ipl)
	}

	for _, r := range w.renderers {
		r.Render(ipl)
	}

	for _, r := range w.postRenderers {
		r.PostRender(ipl)
	}
}
