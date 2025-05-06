package events

import (
	"sync"

	"go.uber.org/zap"
)

type EventBus interface {
	RegisterHandler(handler EventHandler)
	Publish(event Event)
	GetEventLog() []Event
}

type MemoryEventBus struct {
	mu         sync.RWMutex
	handlers   []EventHandler
	eventLog   []Event
	queue      []Event
	processing bool
}

// NewEventBus creates a new event bus
func NewMemoryEventBus() EventBus {
	return &MemoryEventBus{
		handlers: make([]EventHandler, 0),
		eventLog: make([]Event, 0),
		queue:    make([]Event, 0),
	}
}

// RegisterHandler registers a handler for all events
func (eb *MemoryEventBus) RegisterHandler(handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.handlers = append(eb.handlers, handler)
}

// Publish publishes an event to the bus
func (eb *MemoryEventBus) Publish(event Event) {
	eb.mu.Lock()
	eb.queue = append(eb.queue, event)
	eb.mu.Unlock()

	// Start processing if not already processing
	if !eb.processing {
		eb.processQueue()
	}
}

// processQueue processes all events in the queue
func (eb *MemoryEventBus) processQueue() {
	eb.mu.Lock()
	if eb.processing {
		eb.mu.Unlock()
		return
	}
	eb.processing = true
	eb.mu.Unlock()

	defer func() {
		eb.mu.Lock()
		eb.processing = false
		eb.mu.Unlock()
	}()

	for {
		eb.mu.Lock()
		if len(eb.queue) == 0 {
			eb.mu.Unlock()
			return
		}

		// Get the next event
		event := eb.queue[0]
		eb.queue = eb.queue[1:]
		// Log the event
		eb.eventLog = append(eb.eventLog, event)
		eb.mu.Unlock()

		// Get all handlers
		eb.mu.RLock()
		handlers := eb.handlers
		eb.mu.RUnlock()

		// Process the event through all handlers
		for _, handler := range handlers {
			newEvents, err := handler.Handle(event)
			if err != nil {
				zap.L().Error("error handling event", zap.Error(err))
				continue
			}
			if len(newEvents) > 0 {
				eb.mu.Lock()
				eb.queue = append(eb.queue, newEvents...)
				eb.mu.Unlock()
				eb.processQueue()
			}
		}
	}
}

// GetEventLog returns the complete event log
func (eb *MemoryEventBus) GetEventLog() []Event {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	events := make([]Event, len(eb.eventLog))
	copy(events, eb.eventLog)
	return events
}
