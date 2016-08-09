package eventsourcing

// EventBus is an interface defining an event bus for distributing events.
type EventBus interface {
	// HandleEvents start event bus
	HandleEvents()
	AddGlobalHandler(Handler)
	RegisterHandler(Handler)
	RegisterHandlers(interface{})
}

type InternalEventBus struct {
	_store      EventStore
	_dispatcher Dispatcher
}

func NewInternalEventBus(store EventStore) EventBus {
	return &InternalEventBus{
		_store:      store,
		_dispatcher: NewDispatcher(),
	}
}

func (this InternalEventBus) AddGlobalHandler(handler Handler) {
	this._dispatcher.AddGlobalHandler(handler)
}

func (this InternalEventBus) RegisterHandler(handler Handler) {
	this._dispatcher.RegisterHandler(handler)
}

func (this InternalEventBus) RegisterHandlers(source interface{}) {
	this._dispatcher.RegisterHandlers(source)
}

func (this InternalEventBus) HandleEvents() {
	eventChannel, ok := this._store.(EventChannel)
	if !ok {
		return
	}
	eventChan := eventChannel.GetEventChan()
	for {
		event := <-eventChan
		this._dispatcher.Dispatch(event)
	}
}
