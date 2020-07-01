package network

type EventManager interface {
	Queue(event interface{})
	Flush()
}

type System struct {
	EventManager EventManager
}

func CreateSystem(em EventManager) *System {
	return &System{
		EventManager: em,
	}
}

func (s *System) Queue(event interface{}) {
	s.EventManager.Queue(event)
}

func (s *System) PreSimulate(dt float64) {
	s.EventManager.Flush()
}
