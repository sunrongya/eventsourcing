package eventsourcing

import (
	"fmt"
	"reflect"
	"strings"
)

var ErrEventNotRegistered = fmt.Errorf("event not registered")
var ErrEventAlreadyRegiter = fmt.Errorf("event already registered")
var ErrCreateEvent = fmt.Errorf("create event error")

type EventFactory struct {
	_factories map[string]reflect.Type
}

func NewEventFactory() *EventFactory {
	return &EventFactory{
		_factories: make(map[string]reflect.Type),
	}
}

func (this *EventFactory) RegisterEventType(eventType reflect.Type) error {
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}
	if _, ok := this._factories[eventType.String()]; ok {
		return ErrEventAlreadyRegiter
	}
	this._factories[eventType.String()] = eventType

	return nil
}

func (this *EventFactory) RegisterAggregate(aggregates ...Aggregate) {
	for _, aggregate := range aggregates {
		aggregateType := reflect.TypeOf(aggregate)
		for i := 0; i < aggregateType.NumMethod(); i++ {
			method := aggregateType.Method(i)
			if strings.HasPrefix(method.Name, "Handle") {
				this.RegisterEventType(method.Type.In(1))
			}
		}
	}
}

func (this *EventFactory) GetEvent(eventType string) (Event, error) {
	reflectEventType, ok := this._factories[eventType]
	if !ok {
		return nil, ErrEventNotRegistered
	}
	event, ok := reflect.New(reflectEventType).Interface().(Event)
	if !ok {
		return nil, ErrCreateEvent
	}
	return event, nil
}

func (this *EventFactory) EventStringType(event Event) string {
	eventType := reflect.TypeOf(event)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}
	return eventType.String()
}
