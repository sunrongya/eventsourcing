package eventsourcing

import (
	"reflect"
	"strings"
)

type ProcessCommandFun func(Command) []Event
type ARProcessCommandFun func(interface{}) ProcessCommandFun

func NewCommandDispatcher() *CommandDispatcher {
	return &CommandDispatcher{
		_prefix:   "Process",
		_handlers: make(map[reflect.Type]ARProcessCommandFun),
	}
}

type CommandDispatcher struct {
	_prefix   string
	_handlers map[reflect.Type]ARProcessCommandFun
}

func (this *CommandDispatcher) add(commandType reflect.Type, handler ARProcessCommandFun) {
	this._handlers[commandType] = handler
}

func (this *CommandDispatcher) Register(source interface{}) {
	aggregateType := reflect.TypeOf(source)
	for i := 0; i < aggregateType.NumMethod(); i++ {
		method := aggregateType.Method(i)
		if !strings.HasPrefix(method.Name, this._prefix) {
			continue
		}
		commadType := method.Type.In(1)
		handler := func(aggerate interface{}) ProcessCommandFun {
			return func(commad Command) (events []Event) {
				commadValue := reflect.ValueOf(commad)
				values := method.Func.Call([]reflect.Value{
					reflect.ValueOf(aggerate),
					commadValue})
				for _, v := range values {
					events = append(events, v.Interface().([]Event)...)
				}
				return
			}
		}
		this.add(commadType, handler)
	}
}

func (this *CommandDispatcher) Get(command Command) (ARProcessCommandFun, bool) {
	handler, ok := this._handlers[reflect.TypeOf(command)]
	return handler, ok
}
