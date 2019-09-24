package engine

import (
	"testing"
	"time"
)

type mockStepHandler struct {
	t             *testing.T
	updates       int
	previousDt    float64
	previousAlpha float64
	renders       int
	checkOrder    bool
}

func (e *mockStepHandler) Simulate(dt float64) {
	e.updates++
	e.previousDt = dt
}

func (e *mockStepHandler) Render(ipl float64) {
	// one render behind when this assertions triggers
	if e.checkOrder && e.renders != e.updates {
		e.t.Fatalf(
			"invalid call order, Render did not run after Simulate: %d renders %d updates",
			e.renders,
			e.updates,
		)
	}

	e.renders++
	e.previousAlpha = ipl
}

func TestFixedStepSimulation(t *testing.T) {
	t.Run("calls Simulate at least once per Step", func(t *testing.T) {
		msh := &mockStepHandler{t: t}
		fss := fixedStepSimulation{sps: 1}

		fss.Step(msh, msh, 0)
		fss.Step(msh, msh, 1*time.Second)

		got := msh.updates
		want := 1
		if got != want {
			t.Errorf("incorrect number of calls to Simulate: got %d expected %d", got, want)
		}
	})

	t.Run("calls Render independently of Simulate", func(t *testing.T) {
		msh := &mockStepHandler{t: t}
		fss := fixedStepSimulation{sps: 1}

		fss.Step(msh, msh, 0)
		fss.Step(msh, msh, time.Second/2)

		got := msh.renders
		want := 2
		if got != want {
			t.Errorf("incorrect number of calls to Render: got %d expected %d", got, want)
		}
	})

	t.Run("calls Render after Simulate", func(t *testing.T) {
		msh := &mockStepHandler{t: t, checkOrder: true}
		fss := fixedStepSimulation{sps: 1}

		fss.Step(msh, msh, 0)
		fss.Step(msh, msh, time.Second)

		got := msh.renders
		want := 2
		if got != want {
			t.Errorf("incorrect number of calls to Render: got %d expected %d", got, want)
		}
	})

	t.Run("passes a fixed delta time to Simulate", func(t *testing.T) {
		msh := &mockStepHandler{t: t}
		fss := fixedStepSimulation{sps: 20}

		want := time.Second.Seconds() / 20
		for i := 0; i < 20; i++ {
			fss.Step(msh, msh, 0)
			fss.Step(msh, msh, time.Second)
			got := msh.previousDt
			if got != want {
				t.Errorf("incorrect delta time: got %f expected %f", got, want)
			}
		}
	})

	t.Run("calls Simulate extra times if falling behind on Simulations", func(t *testing.T) {
		msh := &mockStepHandler{t: t}
		fss := fixedStepSimulation{sps: 10}

		fss.Step(msh, msh, 0)
		fss.Step(msh, msh, 2*time.Second)
		got := msh.updates
		want := 2
		if got != want {
			t.Errorf("incorrect number of calls to Simulate: got %d expected %d", got, want)
		}
	})

	t.Run("calls Simulate up to a maximum of twenty percent of simulations per second times", func(t *testing.T) {
		msh := &mockStepHandler{t: t}
		fss := fixedStepSimulation{sps: 20}
		fss.Step(msh, msh, 0)
		fss.Step(msh, msh, 251*time.Millisecond)
		got := msh.updates
		want := 4
		if got != want {
			t.Errorf("incorrect number of calls to Simulate: got %d expected %d", got, want)
		}
	})

	t.Run("passes the frame progress as the alpha to Render", func(t *testing.T) {
		alphaTests := []struct {
			sps     uint
			elapsed time.Duration
			want    float64
		}{
			{sps: 20, elapsed: 1 * time.Millisecond, want: 0.02},
			{sps: 20, elapsed: 2 * time.Millisecond, want: 0.04},
			{sps: 20, elapsed: 3 * time.Millisecond, want: 0.06},
			{sps: 20, elapsed: 4 * time.Millisecond, want: 0.08},
			{sps: 20, elapsed: 5 * time.Millisecond, want: 0.10},
			{sps: 20, elapsed: 25 * time.Millisecond, want: 0.50},
			{sps: 20, elapsed: 50 * time.Millisecond, want: 0.00},
			{sps: 20, elapsed: 55 * time.Millisecond, want: 0.10},
			{sps: 20, elapsed: 251 * time.Millisecond, want: 1.02},
			{sps: 1, elapsed: 500 * time.Millisecond, want: 0.50},
		}

		for _, tc := range alphaTests {
			msh := &mockStepHandler{t: t}
			fss := fixedStepSimulation{sps: tc.sps}
			fss.Step(msh, msh, 0)
			fss.Step(msh, msh, tc.elapsed)
			got := msh.previousAlpha
			if got != tc.want {
				t.Errorf("previous alpha value passed to render is incorrect: got %f expected %f", got, tc.want)
			}
		}
	})
}
