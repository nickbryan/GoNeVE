package input

import (
	"reflect"
	"testing"

	"github.com/nickbryan/GoNeVE/pkg/event"
)

type mockSubscriber struct {
	cb event.Callback
}

func (ms *mockSubscriber) Subscribe(cb event.Callback, topics ...event.Topic) {
	ms.cb = cb
}

func TestNew(t *testing.T) {
	t.Run("returns a pointer to a new System", func(t *testing.T) {
		// TODO: assert nil return on all these type of tests
		got := reflect.TypeOf(CreateSystem(&mockSubscriber{}))
		want := reflect.TypeOf(&System{})
		if got != want {
			t.Errorf("got %v expected %v", got, want)
		}
	})
}

type command struct {
	name   string
	key    Key
	state  State
	action Action
	event  event.Topic
	msg    string
	sims   int
}

func TestSystem(t *testing.T) {
	tests := []command{
		{
			name:   "callback is triggered on key press",
			key:    KeyA,
			state:  Press,
			action: KeyPressed,
			event:  KeyPressedEvent,
			msg:    "KeyA pressed",
			sims:   1,
		},
		{
			name:   "callback is triggered on key release with different key",
			key:    KeyB,
			state:  Press,
			action: KeyPressed,
			event:  KeyPressedEvent,
			msg:    "KeyB pressed",
			sims:   1,
		},
		{
			name:   "callback is triggered on key release",
			key:    KeyA,
			state:  Release,
			action: KeyReleased,
			event:  KeyReleasedEvent,
			msg:    "KeyA released",
			sims:   1,
		},
		{
			name:   "callback is triggered on key release with different key",
			key:    KeyC,
			state:  Release,
			action: KeyReleased,
			event:  KeyReleasedEvent,
			msg:    "KeyC released",
			sims:   1,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ms := &mockSubscriber{}
			s := CreateSystem(ms)

			var msg string
			fdt := float64(i)

			s.AddKeyCommands(tc.key, tc.state, KeyCommandExecutorFunc(func(dt float64) {
				msg = tc.msg
				if dt != fdt {
					t.Errorf("delta time value unexpecetd: got %f expectd %f", dt, fdt)
				}
			}))

			ms.cb(KeyEventMessage{
				Action: tc.action,
				Key:    tc.key,
			}, tc.event)

			s.PreSimulate(fdt)

			if msg != tc.msg {
				t.Errorf("message %s does not match expected %s", msg, tc.msg)
			}
		})
	}

	t.Run("callback is triggered over multiple PreSimulate calls when Pressed until Released", func(t *testing.T) {
		ms := &mockSubscriber{}
		s := CreateSystem(ms)

		var msg string
		want := "KeyA pressed"

		s.AddKeyCommands(KeyA, Pressed, KeyCommandExecutorFunc(func(dt float64) {
			msg = "KeyA pressed"
			if dt != 10.5 {
				t.Errorf("delta time value unexpecetd: got %f expectd %f", dt, 10.5)
			}
		}))

		ms.cb(KeyEventMessage{
			Action: KeyPressed,
			Key:    KeyA,
		}, KeyPressedEvent)

		for i := 0; i < 10; i++ {
			s.PreSimulate(10.5)
			if msg != want {
				t.Errorf("message %s does not match expected %s", msg, want)
			}
		}

		msg = "should not change"

		ms.cb(KeyEventMessage{
			Action: KeyReleased,
			Key:    KeyA,
		}, KeyReleasedEvent)

		s.PreSimulate(10.5)
		if msg == want {
			t.Errorf("callback triggered: message %s does not match expected %s", msg, want)
		}
	})
}
