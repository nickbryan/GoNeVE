package main

import (
	"log"
	"sync"
	"time"
)

type Manager struct {
	queue        chan string
	flush        chan struct{}
	queuedEvents []string
}

func (m *Manager) Initialise() {
	go func() {
		for {
			select {
			case topic := <-m.queue:
				m.queuedEvents = append(m.queuedEvents, topic)
			case <-m.flush:
				m.queuedEvents = nil
			}
		}
	}()
}

func (m *Manager) Queue(topic string) {
	m.queue <- topic
}

func (m *Manager) Flush() {
	if len(m.queuedEvents) > 0 {
		start := time.Now()
		events := make([]string, len(m.queuedEvents))
		copy(events, m.queuedEvents)
		m.flush <- struct{}{}
		count := 0
		for _, t := range events {
			log.Println(t)
			count++
		}
		log.Println(count, "------------------------------------------------------------", time.Since(start))
	}
}

func main() {
	m := &Manager{
		queue: make(chan string),
		flush: make(chan struct{}),
	}

	m.Initialise()

	w := sync.WaitGroup{}
	w.Add(2)
	go func() {
		for {
			wg := sync.WaitGroup{}
			wg.Add(2)
			go func() {
				m.Queue("a")
				m.Queue("b")
				m.Queue("c")
				m.Queue("d")
				m.Queue("e")
				m.Queue("f")
				m.Queue("g")
				m.Queue("h")
				m.Flush()
				m.Queue("i")
				m.Queue("j")
				m.Queue("k")
				m.Queue("l")
				wg.Done()
			}()

			go func() {
				m.Queue("m")
				m.Queue("n")
				m.Queue("o")
				m.Flush()
				m.Queue("p")
				m.Queue("q")
				m.Queue("r")
				wg.Done()
			}()
			wg.Wait()
			m.Queue("batch complete")
		}
	}()

	go func() {
		t := time.NewTicker(15 * time.Millisecond)
		for range t.C {
			m.Flush()
		}
	}()
	w.Wait()
}
