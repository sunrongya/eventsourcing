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
	store      *MemEventStore
	dispatcher Dispatcher
}

func NewInternalEventBus(store *MemEventStore) EventBus {
	return &InternalEventBus{
		store:      store,
		dispatcher: NewDispatcher(),
	}
}

func (this InternalEventBus) AddGlobalHandler(handler Handler) {
	this.dispatcher.AddGlobalHandler(handler)
}

func (this InternalEventBus) RegisterHandler(handler Handler) {
	this.dispatcher.RegisterHandler(handler)
}

func (this InternalEventBus) RegisterHandlers(source interface{}) {
	this.dispatcher.RegisterHandlers(source)
}

func (this InternalEventBus) HandleEvents() {
	eventChan := this.store.GetEventChan()
	for {
		event := <-eventChan
		this.dispatcher.Dispatch(event)
	}
}
