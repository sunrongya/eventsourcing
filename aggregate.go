package eventsourcing

import (
	"fmt"
)

// Common interface for all event-sourced aggregates
type Aggregate interface {
	Guider

	Version() int
	SetVersion(int)

	// apply a list of events to restore actual state (core event sourcing)
	// ApplyEvents([]Event)

	// Process a command according to own actual state (eg. debit account checks account.balance)
	// Produce proper state-changing events
	// ProcessCommand(Command) []Event
}

// base implementation for all aggregates - with GUID and Version
type BaseAggregate struct {
	WithGuid
	_version int
}

func (this *BaseAggregate) Version() int {
	return this._version
}

func (this *BaseAggregate) SetVersion(version int) {
	this._version = version
}

// restores given empty aggregate from a state stored in event store
func RestoreAggregate(guid Guid, a Aggregate, store EventStore) {
	events, _ := store.Find(guid)
	for _, event := range events {
		handler, ok := globalEventDispatcher().Get(a, event)
		if ok {
			handler(a)(event)
		} else {
			panic(fmt.Errorf("Unknown event %#v", event))
		}
	}
	//a.ApplyEvents(events)
	a.SetVersion(len(events))
	a.SetGuid(guid)
}
