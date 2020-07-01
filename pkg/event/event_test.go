package event

import (
	"reflect"
	"sync"
	"testing"
)

// TODO: refactor message to be an interface that has a type function that should be set like engo?
// TODO: improve concurrent call performance of the benchmark, maybe use channels and a go routine to manage them?

func TestNewManager(t *testing.T) {
	t.Run("returns a pointer to a new Manager", func(t *testing.T) {
		got := reflect.TypeOf(NewManager())
		want := reflect.TypeOf(&Manager{})
		if got != want {
			t.Errorf("got %v expected %v", got, want)
		}
	})
}

func TestManager(t *testing.T) {
	assertCalls := func(t *testing.T, want int, assert func(c *int)) {
		t.Helper()
		calls := 0
		assert(&calls)
		if calls != want {
			t.Errorf("callback called incorrect number of times: got %d want %d", calls, want)
		}
	}

	t.Run("subscribes a callback to a topic", func(t *testing.T) {
		assertCalls(t, 1, func(c *int) {
			manager := NewManager()
			topic := Topic("topic")

			manager.Subscribe(func(msg interface{}, _ Topic) {
				*c += msg.(int)
			}, topic)

			manager.Publish(1, topic)
		})
	})

	t.Run("subscribes a callback to multiple topics", func(t *testing.T) {
		assertCalls(t, 2, func(c *int) {
			manager := NewManager()
			topicA := Topic("topicA")
			topicB := Topic("topicB")

			manager.Subscribe(func(msg interface{}, _ Topic) {
				*c += msg.(int)
			}, topicA, topicB)

			manager.Publish(1, topicA)
			manager.Publish(1, topicB)
		})
	})

	t.Run("subscribes multiple callbacks to a topic", func(t *testing.T) {
		assertCalls(t, 2, func(c *int) {
			manager := NewManager()
			topic := Topic("topic")

			manager.Subscribe(func(msg interface{}, _ Topic) {
				*c += msg.(int)
			}, topic)
			manager.Subscribe(func(msg interface{}, _ Topic) {
				*c += msg.(int)
			}, topic)

			manager.Publish(1, topic)
		})
	})

	t.Run("publishes multiple topics", func(t *testing.T) {
		assertCalls(t, 2, func(c *int) {
			manager := NewManager()
			topicA := Topic("topicA")
			topicB := Topic("topicB")

			manager.Subscribe(func(msg interface{}, _ Topic) {
				*c += msg.(int)
			}, topicA, topicB)
			manager.Publish(1, topicA, topicB)
		})
	})

	t.Run("publishes a message to subscribers of a topic", func(t *testing.T) {
		assertCalls(t, 1, func(c *int) {
			manager := NewManager()
			topic := Topic("topic")

			manager.Subscribe(func(msg interface{}, _ Topic) {
				if _, ok := msg.(int); !ok {
					t.Fatalf("msg is of incorrect type: got %T want %T", msg, 1)
				}
				*c += msg.(int)
			}, topic)

			manager.Publish(1, topic)
		})
	})

	t.Run("sends the topic that fired the callback to the callback", func(t *testing.T) {
		manager := NewManager()
		topic := Topic("some.topic")

		manager.Subscribe(func(msg interface{}, tpc Topic) {
			if tpc != topic {
				t.Fatalf("recieved topic was not expected: got %s want %s", tpc, topic)
			}

		}, topic)

		manager.Publish(true, topic)
	})

	t.Run("queues events to be flushed later", func(t *testing.T) {
		assertCalls(t, 3, func(c *int) {
			manager := NewManager()
			topicA := Topic("topicA")
			topicB := Topic("topicB")

			manager.Subscribe(func(msg interface{}, _ Topic) {
				*c += msg.(int)
			}, topicA, topicB)

			manager.Queue(1, topicA)
			manager.Queue(1, topicA)
			manager.Queue(1, topicB)

			manager.Flush()
		})
	})
}

func BenchmarkManager_Queue(b *testing.B) {
	manager := NewManager()

	event := struct {
		test int
	}{test: 1}

	for n := 0; n < b.N; n++ {
		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go manager.Queue(event, "topic")
			wg.Done()
		}
		wg.Wait()
	}
}
