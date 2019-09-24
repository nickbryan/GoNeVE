package input

import (
	"github.com/nickbryan/GoNeVE/pkg/engine"
	"github.com/nickbryan/GoNeVE/pkg/event"
)

type keyState int

const (
	press keyState = iota
	release
	off
)

// System is responsible for executing all input commands.
//
// System will subscribe to the KeyPressedEvent and KeyReleasedEvent topics of the supplied engine.Subscriber
// upon creation. When the System is notified of a KeyEventMessage it will call the KeyCommandExecutor registered
// with the specified key and action.
//
// All KeyCommandExecutor's are called within the PreSimulate method to ensure that they only get called once per simulation.
type System struct {
	keys              map[Key]keyState
	keyPressListeners []keyPressListener
}


// CreateSystem creates a new System and subscribes the keyCallback to the KeyPressedEvent and KeyReleasedEvent
// via the passed in Subscriber.
func CreateSystem(subscriber engine.Subscriber) *System {
	s := &System{keys: make(map[Key]keyState)}
	subscriber.Subscribe(s.keyCallback, KeyPressedEvent, KeyReleasedEvent)
	return s
}

// AddKeyCommand will register the given commands, key and state within the Manager. The commands will
// be called once per simulation if the Manager has been notified of the relevant key and state changes.
func (s *System) AddKeyCommands(key Key, state State, commands ...KeyCommandExecutor) {
	for _, l := range s.keyPressListeners {
		if l.key == key && l.state == state {
			l.commands = append(l.commands, commands...)
			return
		}
	}

	s.keys[key] = off

	s.keyPressListeners = append(s.keyPressListeners, keyPressListener{
		key:      key,
		state:    state,
		commands: commands,
	})
}

// PreSimulate is responsible for triggering the registered KeyCommandExecutor's and should be called once per update.
//
// The supplied dt (delta time) will be passed into all KeyCommandExecutor's.
func (s *System) PreSimulate(dt float64) {
	for _, l := range s.keyPressListeners {
		key := s.keys[l.key]

		if key == off {
			continue
		}

		if l.state == Pressed && key == press {
			for _, c := range l.commands {
				c.Execute(dt)
			}
		}

		if l.state == Press && key == press {
			for _, c := range l.commands {
				c.Execute(dt)
			}

			s.keys[l.key] = off
		}

		if l.state == Release && key == release {
			for _, c := range l.commands {
				c.Execute(dt)
			}

			s.keys[l.key] = off
		}
	}
}

func (s *System) keyCallback(msg interface{}, _ event.Topic) {
	if msg, ok := msg.(KeyEventMessage); ok {
		if msg.Action == KeyPressed {
			s.keys[msg.Key] = press
		}

		if msg.Action == KeyReleased {
			s.keys[msg.Key] = release
		}
	}
}
