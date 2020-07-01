package network

import (
	"reflect"
	"testing"

	"github.com/nickbryan/GoNeVE/pkg/engine"
)

type eventStub struct {
	called bool
}

type eventManagerStub struct {
	queue []interface{}
}

func (em *eventManagerStub) Queue(event interface{}) {
	em.queue = append(em.queue, event)
}

func (em *eventManagerStub) Flush() {
	for _, e := range em.queue {
		e.(*eventStub).called = true
	}
}

func TestCreateSystem(t *testing.T) {
	t.Run("returns a pointer to a new System", func(t *testing.T) {
		got := reflect.TypeOf(CreateSystem(&eventManagerStub{}))
		want := reflect.TypeOf(&System{})
		if got != want {
			t.Errorf("got %v expected %v", got, want)
		}
	})
}

func TestSystem(t *testing.T) {
	t.Run("queues up an event to be dispatched later", func(t *testing.T) {
		em := &eventManagerStub{}
		system := CreateSystem(em)
		system.Queue(&eventStub{})
		got := len(em.queue)
		want := 1
		if got != want {
			t.Errorf("event was not added to Queuer: got %d expected %d", got, want)
		}
	})

	t.Run("can be used as a engine.PreSimulator", func(t *testing.T) {
		preSimulator := reflect.TypeOf((*engine.PreSimulator)(nil)).Elem()
		system := CreateSystem(&eventManagerStub{})
		if !reflect.TypeOf(system).Implements(preSimulator) {
			t.Errorf("%T is not a %v", system, preSimulator)
		}
	})

	t.Run("flushes queued events at the start of a simulation", func(t *testing.T) {
		event := &eventStub{}
		em := &eventManagerStub{queue: []interface{}{event}}
		system := CreateSystem(em)
		system.PreSimulate(10.0)
		if !event.called {
			t.Errorf("event was never flushed by PreSimulate")
		}
	})
}
