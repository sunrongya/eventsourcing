package eventsourcing

import (
	"errors"
	"fmt"
	"math"
)

// Event store's common interface
type EventStore interface {
	// find all events for given ID (aggregate).
	// returns event list as well as aggregate version
	Find(guid Guid) (events []Event, version int)

	// Update an aggregate with new events. If the version specified
	// does not match with the version in the Event Store, an error is returned
	Update(guid Guid, version int, events []Event) error

	// Get events from Event Store.
	// Supports pagination with use of offset and batchsize.
	GetEvents(guid Guid, offset int, batchSize int) []Event
}

//in-memory event store. Uses slice for 'complete events catalogue'
// and a map for 'per aggregate' events
type MemEventStore struct {
	EventChannel
	_store  map[Guid][]Event
	_events []Event
}

// @see EventStore.GetEvents
func (this *MemEventStore) Find(guid Guid) ([]Event, int) {
	events := this._store[guid]
	return events, len(events)
}

// @see EventStore.Update
func (this *MemEventStore) Update(guid Guid, version int, events []Event) error {
	changes, ok := this._store[guid]
	if !ok {
		// initialize if not exists
		changes = []Event{}
	}
	if len(changes) == version {
		for _, event := range events {
			event.SetGuid(guid)
		}
		this.AppendEvents(events)
		this._store[guid] = append(changes, events...)
	} else {
		return errors.New(
			fmt.Sprintf("Optimistic locking exeption - client has version %v, but store %v", version, len(changes)))
	}
	return nil
}

// @see EventStore.GetEvents
func (this *MemEventStore) GetEvents(guid Guid, offset int, batchSize int) []Event {
	until := int(math.Min(float64(offset+batchSize), float64(len(this._events))))
	return this._events[offset:until]
}

// initializer for event store
func NewInMemStore() *MemEventStore {
	return &MemEventStore{
		EventChannel: NewEventChannel(),
		_store:       map[Guid][]Event{},
		_events:      make([]Event, 0),
	}
}
