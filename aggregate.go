package eventsourcing

// Common interface for all event-sourced aggregates
type Aggregate interface {
	Guider

	Version() int
	// apply a list of events to restore actual state (core event sourcing)
	ApplyEvents([]Event)

	// Process a command according to own actual state (eg. debit account checks account.balance)
	// Produce proper state-changing events
	ProcessCommand(Command) []Event
}

// base implementation for all aggregates - with GUID and Version
type BaseAggregate struct {
	WithGuid
	version int
}

func (b *BaseAggregate) Version() int {
	return b.version
}

func (b *BaseAggregate) SetVersion(version int) {
	b.version = version
}

// restores given empty aggregate from a state stored in event store
func RestoreAggregate(guid Guid, a Aggregate, store EventStore) {
	events, _ := store.Find(guid)
	a.ApplyEvents(events)
	a.SetGuid(guid)
}
