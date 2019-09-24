package event

import "sync"

// Topic should be a unique name that is used to reference a set of events published on a given channel.
type Topic string

// Callback is the function that will be called when a message is published to a topic.
type Callback func(message interface{}, topic Topic)

// Manager provides a publish/subscribe based event handler that can be used for system communication.
// A Queue and Flush mechanism is also provided to allow events to be queued and flushed when required.
type Manager struct {
	subscriptions map[Topic][]Callback
	queue         map[Topic][]interface{}
	queueMux      sync.Mutex
}

// NewManager will create and initialise a new *Manager.
func NewManager() *Manager {
	return &Manager{
		subscriptions: make(map[Topic][]Callback),
		queue:         make(map[Topic][]interface{}),
	}
}

// Subscribe a Callback to one ore more Topics.
func (m *Manager) Subscribe(cb Callback, topics ...Topic) {
	for _, t := range topics {
		m.subscriptions[t] = append(m.subscriptions[t], cb)
	}
}

// Publish a message to one or more Topics.
func (m *Manager) Publish(msg interface{}, topics ...Topic) {
	for _, t := range topics {
		for _, cb := range m.subscriptions[t] {
			cb(msg, t)
		}
	}
}

// Queue a message to be flushed later on.
func (m *Manager) Queue(msg interface{}, topic Topic) {
	m.queueMux.Lock()
	m.queue[topic] = append(m.queue[topic], msg)
	m.queueMux.Unlock()
}

// Flush all queued messages at once.
func (m *Manager) Flush() {
	for t, msgs := range m.queue {
		for _, msg := range msgs {
			for _, cb := range m.subscriptions[t] {
				cb(msg, t)
			}
		}
	}
}
