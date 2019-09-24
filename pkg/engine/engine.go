package engine

import "github.com/nickbryan/GoNeVE/pkg/event"

// Publisher is the interface that wraps the Publish method.
//
// Publish sends the specified message to the callbacks subscribed to the specified topics.
type Publisher interface {
	Publish(msg interface{}, topics ...event.Topic)
}

// Subscriber is the interface that wraps the Subscribe method.
//
// Subscribe subscribed the specified callback to the specified topics.
type Subscriber interface {
	Subscribe(cb event.Callback, topics ...event.Topic)
}

// PublishSubscriber is a basic publish/subscribe event system. It is used
// for communication between system within the engine.
type PublishSubscriber interface {
	Publisher
	Subscriber
}

// Engine is responsible for handling the main game loop and managing the World.
type Engine struct {
	simStepper SimulationStepper
	world      *World
	running    bool
}

// New creates a new Engine instance and applies the supplied Options to the Engine.
func New(opts ...Option) *Engine {
	options := defaultEngineOptions()

	for _, opt := range opts {
		opt.apply(&options)
	}

	e := &Engine{
		simStepper: options.simStepper,
	}

	e.world = NewWorld()

	return e
}

// Run will call the SimulationStepper until the Engine is told to Stop.
func (e *Engine) Run() {
	e.running = true

	for e.running {
		e.simStepper.Step(e.world, e.world, 0)
	}
}

// Stop will stop the main game loop.
func (e *Engine) Stop() {
	e.running = false
}
