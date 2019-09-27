package engine

import (
	"reflect"
	"testing"
)

func TestNewWorld(t *testing.T) {
	t.Run("returns a pointer to a new World", func(t *testing.T) {
		got := reflect.TypeOf(NewWorld())
		want := reflect.TypeOf(&World{})
		if got != want {
			t.Errorf("got %v expected %v", got, want)
		}
	})
}

type callAssertion struct {
	callType string
	dt       float64
}

type stub struct {
	calls []callAssertion
}

func (s *stub) Simulate(dt float64) {
	s.calls = append(s.calls, callAssertion{
		callType: "simulate",
		dt:       dt,
	})
}

func (s *stub) PreSimulate(dt float64) {
	s.calls = append(s.calls, callAssertion{
		callType: "pre_simulate",
		dt:       dt,
	})
}

func (s *stub) PostSimulate(dt float64) {
	s.calls = append(s.calls, callAssertion{
		callType: "post_simulate",
		dt:       dt,
	})
}

func (s *stub) Render(ipl float64) {
	s.calls = append(s.calls, callAssertion{
		callType: "render",
		dt:       ipl,
	})
}

func (s *stub) PreRender(ipl float64) {
	s.calls = append(s.calls, callAssertion{
		callType: "pre_render",
		dt:       ipl,
	})
}

func (s *stub) PostRender(ipl float64) {
	s.calls = append(s.calls, callAssertion{
		callType: "post_render",
		dt:       ipl,
	})
}

func TestWorld(t *testing.T) {
	assertCalls := func(t *testing.T, got, want []callAssertion) {
		t.Helper()
		if len(got) != len(want) {
			t.Errorf("number of calls don't match: got %d want %d", len(got), len(want))
		}
		for i, v := range got {
			if v != want[i] {
				t.Errorf("expected calls do not match: got %v want %v", got, want)
			}
		}
	}
	t.Run("can add a Simulator", func(t *testing.T) {
		s := &stub{}
		world := NewWorld()
		world.AddSimulator(s)
		world.Simulate(10.0)
		assertCalls(t, s.calls, []callAssertion{
			{callType: "simulate", dt: 10.0},
		})
	})

	t.Run("can add multiple Simulators", func(t *testing.T) {
		s := &stub{}
		world := NewWorld()
		world.AddSimulator(s)
		world.AddSimulator(s)
		world.Simulate(10.0)
		assertCalls(t, s.calls, []callAssertion{
			{callType: "simulate", dt: 10.0},
			{callType: "simulate", dt: 10.0},
		})
	})

	t.Run("can add a PreSimulator that is called before any Simulators", func(t *testing.T) {
		s := &stub{}
		world := NewWorld()
		world.AddSimulator(s)
		world.AddPreSimulator(s)
		world.Simulate(10.0)
		assertCalls(t, s.calls, []callAssertion{
			{callType: "pre_simulate", dt: 10.0},
			{callType: "simulate", dt: 10.0},
		})
	})

	t.Run("can add multiple PreSimulators", func(t *testing.T) {
		s := &stub{}
		world := NewWorld()
		world.AddSimulator(s)
		world.AddPreSimulator(s)
		world.AddPreSimulator(s)
		world.Simulate(10.0)
		assertCalls(t, s.calls, []callAssertion{
			{callType: "pre_simulate", dt: 10.0},
			{callType: "pre_simulate", dt: 10.0},
			{callType: "simulate", dt: 10.0},
		})
	})

	t.Run("can add a PostSimulator that is called after any Simulators", func(t *testing.T) {
		s := &stub{}
		world := NewWorld()
		world.AddSimulator(s)
		world.AddPostSimulator(s)
		world.Simulate(10.0)
		assertCalls(t, s.calls, []callAssertion{
			{callType: "simulate", dt: 10.0},
			{callType: "post_simulate", dt: 10.0},
		})
	})

	t.Run("can add multiple PostSimulators", func(t *testing.T) {
		s := &stub{}
		world := NewWorld()
		world.AddSimulator(s)
		world.AddPostSimulator(s)
		world.AddPostSimulator(s)
		world.Simulate(10.0)
		assertCalls(t, s.calls, []callAssertion{
			{callType: "simulate", dt: 10.0},
			{callType: "post_simulate", dt: 10.0},
			{callType: "post_simulate", dt: 10.0},
		})
	})

	t.Run("passes delta time to all Simulators", func(t *testing.T) {
		s := &stub{}
		world := NewWorld()
		world.AddPreSimulator(s)
		world.AddSimulator(s)
		world.AddPostSimulator(s)
		world.Simulate(5)
		assertCalls(t, s.calls, []callAssertion{
			{callType: "pre_simulate", dt: 5.0},
			{callType: "simulate", dt: 5.0},
			{callType: "post_simulate", dt: 5.0},
		})
	})

	t.Run("can add a Renderer", func(t *testing.T) {
		s := &stub{}
		world := NewWorld()
		world.AddRenderer(s)
		world.Render(0.5)
		assertCalls(t, s.calls, []callAssertion{
			{callType: "render", dt: 0.5},
		})
	})

	t.Run("can add multiple Renderers", func(t *testing.T) {
		s := &stub{}
		world := NewWorld()
		world.AddRenderer(s)
		world.AddRenderer(s)
		world.Render(0.01)
		assertCalls(t, s.calls, []callAssertion{
			{callType: "render", dt: 0.01},
			{callType: "render", dt: 0.01},
		})
	})

	t.Run("can add a PreRenderer that is called before any Renderers", func(t *testing.T) {
		s := &stub{}
		world := NewWorld()
		world.AddRenderer(s)
		world.AddPreRenderer(s)
		world.Render(0.01)
		assertCalls(t, s.calls, []callAssertion{
			{callType: "pre_render", dt: 0.01},
			{callType: "render", dt: 0.01},
		})
	})

	t.Run("can add multiple PreRenderers", func(t *testing.T) {
		s := &stub{}
		world := NewWorld()
		world.AddRenderer(s)
		world.AddPreRenderer(s)
		world.AddPreRenderer(s)
		world.Render(0.01)
		assertCalls(t, s.calls, []callAssertion{
			{callType: "pre_render", dt: 0.01},
			{callType: "pre_render", dt: 0.01},
			{callType: "render", dt: 0.01},
		})
	})

	t.Run("can add a PostRenderer that is called after any Renderers", func(t *testing.T) {
		s := &stub{}
		world := NewWorld()
		world.AddRenderer(s)
		world.AddPostRenderer(s)
		world.Render(0.01)
		assertCalls(t, s.calls, []callAssertion{
			{callType: "render", dt: 0.01},
			{callType: "post_render", dt: 0.01},
		})
	})

	t.Run("can add multiple PostRenderers", func(t *testing.T) {
		s := &stub{}
		world := NewWorld()
		world.AddRenderer(s)
		world.AddPostRenderer(s)
		world.AddPostRenderer(s)
		world.Render(0.01)
		assertCalls(t, s.calls, []callAssertion{
			{callType: "render", dt: 0.01},
			{callType: "post_render", dt: 0.01},
			{callType: "post_render", dt: 0.01},
		})
	})
}
