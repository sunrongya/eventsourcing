package eventsourcing

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type Event1 struct {
	WithGuid
	Name string
}

type Event2 struct {
	WithGuid
	Name string
}

type Event3 struct {
	WithGuid
	Name string
}

type Aggregate1 struct {
	BaseAggregate
}

func (this *Aggregate1) HandleEvent1Event(event *Event1) {
}

func (this *Aggregate1) dsddEvent(event *Event2) {
}

func (this *Aggregate1) HandleEvent3(event *Event3) {
}

func TestEventFactory(t *testing.T) {
	tests := []struct {
		event   Event
		isExist bool
	}{
		{
			new(Event1),
			true,
		},
		{
			new(Event2),
			false,
		},
		{
			new(Event3),
			true,
		},
	}

	factory := NewEventFactory()
	factory.RegisterAggregate(new(Aggregate1))

	want := map[string]reflect.Type{
		reflect.TypeOf(Event1{}).String(): reflect.TypeOf((*Event1)(nil)).Elem(),
		reflect.TypeOf(Event3{}).String(): reflect.TypeOf((*Event3)(nil)).Elem(),
	}

	assert.Equal(t, want, factory._factories)

	for _, v := range tests {
		_, err := factory.GetEvent(factory.EventStringType(v.event))
		assert.Equal(t, v.isExist, err == nil)
	}
}
