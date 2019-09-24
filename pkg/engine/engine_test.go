package engine

import (
	"reflect"
	"testing"
	"time"
)

type mockStepper struct {
	calls  int
	stopAt int
	e      *Engine
}

func (ms *mockStepper) Step(s Simulator, r Renderer, elapsed time.Duration) {
	if ms.calls == ms.stopAt {
		ms.e.Stop()
		return
	}
	ms.calls++
}

func TestNew(t *testing.T) {
	t.Run("returns a pointer to a new Engine", func(t *testing.T) {
		got := reflect.TypeOf(New())
		want := reflect.TypeOf(&Engine{})
		if got != want {
			t.Errorf("got %v expected %v", got, want)
		}
	})

	t.Run("sets the SimulationStepper from the passed in configuration", func(t *testing.T) {
		ms := &mockStepper{}
		e := New(WithSimulationStepper(ms))
		if _, ok := e.simStepper.(*mockStepper); !ok {
			t.Errorf("%v is not a %T", e.simStepper, ms)
		}
	})

	t.Run("sets SimulationStepper to the default fixedStepSimulator if non passed in the config", func(t *testing.T) {
		e := New()
		if _, ok := e.simStepper.(*fixedStepSimulation); !ok {
			t.Errorf("%v is not a %T", e.simStepper, &fixedStepSimulation{})
		}
	})

	t.Run("sets the desired simulations per second on the default SimulationStepper", func(t *testing.T) {
		var want uint = 60
		e := New(WithSimulationsPerSecond(want))
		ss := e.simStepper.(*fixedStepSimulation)
		if ss.sps != want {
			t.Errorf("got %d expected %d", ss.sps, want)
		}
	})
}

func TestRun(t *testing.T) {
	t.Run("calls Step on the SimulationStepper", func(t *testing.T) {
		ms := &mockStepper{stopAt: 1}
		e := New(WithSimulationStepper(ms))
		ms.e = e
		e.Run()
		want := 1
		if ms.calls != want {
			t.Errorf("incorrect number of calls to Run: got %d expected %d", ms.calls, want)
		}
	})

	t.Run("continues to call Step on the SimulationStepper until stopped", func(t *testing.T) {
		ms := &mockStepper{stopAt: 5}
		e := New(WithSimulationStepper(ms))
		ms.e = e
		e.Run()
		want := 5
		if ms.calls != want {
			t.Errorf("incorrect number of calls to Run: got %d expected %d", ms.calls, want)
		}
	})

	// TODO: test world
	// TODO: test elapsed
	// TODO: prevent from being able to run 2 loops
}
