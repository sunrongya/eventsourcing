package eventsourcing

import (
	"reflect"
	"strings"
)

type ProcessCommandFun func(Command) []Event
type ARProcessCommandFun func(interface{}) ProcessCommandFun

func NewCommandDispatcher() *CommandDispatcher {
	return &CommandDispatcher{
		prefix:   "Process",
		handlers: make(map[reflect.Type]ARProcessCommandFun),
	}
}

type CommandDispatcher struct {
	prefix   string
	handlers map[reflect.Type]ARProcessCommandFun
}

func (d *CommandDispatcher) add(commandType reflect.Type, handler ARProcessCommandFun) {
	d.handlers[commandType] = handler
}

func (d *CommandDispatcher) Register(source interface{}) {
	aggregateType := reflect.TypeOf(source)
	for i := 0; i < aggregateType.NumMethod(); i++ {
		method := aggregateType.Method(i)
		if !strings.HasPrefix(method.Name, d.prefix) {
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
				    events = append(events, v.Interface().([]Event)... ) 
				}
				return
			}
		}
		d.add(commadType, handler)
	}
}

func (d *CommandDispatcher) Get(command Command) (ARProcessCommandFun, bool) {
	handler, ok := d.handlers[reflect.TypeOf(command)]
	return handler, ok
}
