package eventsourcing

import (
	"reflect"
	"strings"
)

// 为了可读性，用反射来处理EventHandle分发
// 为了性能，并发操作由使用者控制

var _aggerateEventDispatcher *EventDispatcher

func init() {
	_aggerateEventDispatcher = NewEventDispatcher()
}

func aggerateEventRegister(aggerate interface{}) {
	_aggerateEventDispatcher.Register(aggerate)
}

func globalEventDispatcher() *EventDispatcher {
	return _aggerateEventDispatcher
}

type ApplyEventFun func(Event)
type ARApplyEventFun func(interface{}) ApplyEventFun

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		_prefix:   "Handle",
		_handlers: make(map[reflect.Type]map[reflect.Type]ARApplyEventFun),
	}
}

type EventDispatcher struct {
	_prefix   string
	_handlers map[reflect.Type]map[reflect.Type]ARApplyEventFun
}

func (this *EventDispatcher) add(aggregateType, eventType reflect.Type, handler ARApplyEventFun) {
	if _, ok := this._handlers[aggregateType]; !ok {
		this._handlers[aggregateType] = make(map[reflect.Type]ARApplyEventFun)
	}
	this._handlers[aggregateType][eventType] = handler
}

func (this *EventDispatcher) Register(source interface{}) {
	aggregateType := reflect.TypeOf(source)
	for i := 0; i < aggregateType.NumMethod(); i++ {
		method := aggregateType.Method(i)
		if !strings.HasPrefix(method.Name, this._prefix) {
			continue
		}
		eventType := method.Type.In(1)
		handler := func(aggregate interface{}) ApplyEventFun {
			return func(event Event) {
				eventValue := reflect.ValueOf(event)
				method.Func.Call([]reflect.Value{
					reflect.ValueOf(aggregate),
					eventValue})
			}
		}
		this.add(aggregateType, eventType, handler)
	}
}

func (this *EventDispatcher) IsRegistered(aggregate interface{}) bool {
	_, ok := this._handlers[reflect.TypeOf(aggregate)]
	return ok
}

func (this *EventDispatcher) Get(aggregate interface{}, event Event) (ARApplyEventFun, bool) {
	if _, ok := this._handlers[reflect.TypeOf(aggregate)]; ok {
		handler, ok := this._handlers[reflect.TypeOf(aggregate)][reflect.TypeOf(event)]
		return handler, ok
	}

	return nil, false
}
