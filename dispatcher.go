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
	_messageTypes   map[string]reflect.Type
	_handlers       Handlers
	_globalHandlers GlobalHandlers
}

func (this *InternalDispatcher) AddGlobalHandler(handler Handler) {
	this._globalHandlers = append(this._globalHandlers, handler)
}

func (this *InternalDispatcher) addHandler(messageType reflect.Type, handler Handler) {
	_, ok := this._handlers[messageType]
	if !ok {
		this._handlers[messageType] = make([]Handler, 0)
	}
	_, exists := this._messageTypes[messageType.Name()]
	if !exists {
		this._messageTypes[messageType.Name()] = messageType
	}
	this._handlers[messageType] = append(this._handlers[messageType], handler)
}

func (this *InternalDispatcher) RegisterHandler(handler Handler) {
	eventType := reflect.TypeOf(handler).In(1)
	this.addHandler(eventType, handler)
}

func (this *InternalDispatcher) RegisterHandlers(source interface{}) {
	productType := reflect.TypeOf(source)
	for i := 0; i < productType.NumMethod(); i++ {
		method := productType.Method(i)
		if !strings.HasPrefix(method.Name, "Handle") {
			continue
		}
		eventType := method.Type.In(1)
		handler := func(event Event) {
			eventValue := reflect.ValueOf(event)
			method.Func.Call([]reflect.Value{
				reflect.ValueOf(source),
				eventValue})
		}
		this.addHandler(eventType, handler)
	}
}

func (this *InternalDispatcher) Dispatch(message Event) {
	eventType := reflect.TypeOf(message)
	if val, ok := this._handlers[eventType]; ok {
		for _, handler := range val {
			handler(message)
		}
	}
	for _, handler := range this._globalHandlers {
		handler(message)
	}
}
