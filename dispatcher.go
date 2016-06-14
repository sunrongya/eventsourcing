package eventsourcing

import (
	"reflect"
	"strings"
)

type Handler func(Event)
type Handlers map[reflect.Type][]Handler
type GlobalHandlers []Handler

type Dispatcher interface {
	AddGlobalHandler(Handler)
	RegisterHandler(Handler)
	RegisterHandlers(interface{})
	Dispatch(Event)
}

func NewDispatcher() Dispatcher {
	return &InternalDispatcher{
		make(map[string]reflect.Type),
		make(Handlers),
		make(GlobalHandlers, 0)}
}

type InternalDispatcher struct {
	MessageTypes   map[string]reflect.Type
	handlers       Handlers
	globalHandlers GlobalHandlers
}

func (d *InternalDispatcher) AddGlobalHandler(handler Handler) {
	d.globalHandlers = append(d.globalHandlers, handler)
}

func (d *InternalDispatcher) addHandler(messageType reflect.Type, handler Handler) {
	_, ok := d.handlers[messageType]
	if !ok {
		d.handlers[messageType] = make([]Handler, 0)
	}
	_, exists := d.MessageTypes[messageType.Name()]
	if !exists {
		d.MessageTypes[messageType.Name()] = messageType
	}
	d.handlers[messageType] = append(d.handlers[messageType], handler)
}

func (d *InternalDispatcher) RegisterHandler(handler Handler) {
	eventType := reflect.TypeOf(handler).In(1)
	d.addHandler(eventType, handler)
}

func (d *InternalDispatcher) RegisterHandlers(source interface{}) {

	productType := reflect.TypeOf(source)
	numMethods := productType.NumMethod()

	for i := 0; i < numMethods; i++ {

		method := productType.Method(i)

		if strings.HasPrefix(method.Name, "Handle") {
			eventType := method.Type.In(1)
			handler := func(event Event) {
				eventValue := reflect.ValueOf(event)
				method.Func.Call([]reflect.Value{
					reflect.ValueOf(source),
					eventValue})
			}
			d.addHandler(eventType, handler)
		}

	}

}

func (d *InternalDispatcher) Dispatch(message Event) {
	eventType := reflect.TypeOf(message)
	if val, ok := d.handlers[eventType]; ok {
		for _, handler := range val {
			handler(message)
		}
	}
	for _, handler := range d.globalHandlers {
		handler(message)
	}
}
