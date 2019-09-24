package engine

type (
	engineOptions struct {
		simStepper SimulationStepper
	}

	// Option defines a single configuration option.
	Option interface {
		apply(*engineOptions)
	}

	engineOptionFunc func(*engineOptions)
)

func (eof engineOptionFunc) apply(eo *engineOptions) {
	eof(eo)
}

func defaultEngineOptions() engineOptions {
	return engineOptions{
		simStepper: &fixedStepSimulation{
			sps: defaultSimulationsPerSecond,
		},
	}
}

// WithSimulationStepper sets the SimulationStepper in the configuration.
func WithSimulationStepper(ss SimulationStepper) Option {
	return engineOptionFunc(func(eo *engineOptions) {
		eo.simStepper = ss
	})
}

// WithSimulationsPerSecond allows the consumer to set the desired SimulationsPerSecond for the default SimulationStepper.
func WithSimulationsPerSecond(sps uint) Option {
	return engineOptionFunc(func(eo *engineOptions) {
		if ss, ok := eo.simStepper.(*fixedStepSimulation); ok {
			ss.sps = sps
		}
	})
}
