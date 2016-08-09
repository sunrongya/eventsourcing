package eventsourcing

type EventChannel interface {
	GetEventChan() <-chan Event
	AppendEvents(events []Event)
}

type eventChannel struct {
	_eventChan chan Event
}

func NewEventChannel() EventChannel {
	return &eventChannel{
		_eventChan: make(chan Event, 100),
	}
}

// Get persisted events channel -
// channel notifies of any change persisted int the event store
func (this *eventChannel) GetEventChan() <-chan Event {
	return this._eventChan
}

// Add events to the store and send them down the channel
func (this *eventChannel) AppendEvents(events []Event) {
	for _, e := range events {
		this._eventChan <- e
	}
}
