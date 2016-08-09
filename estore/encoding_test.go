package estore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	es "github.com/sunrongya/eventsourcing"
)

type mockEvent1 struct {
	es.WithGuid
	Name string
}

type mockEvent2 struct {
	es.WithGuid
	Name string
}

type mockNotRegisterEvent struct {
	es.WithGuid
	Name string
}

type mockAggregate struct {
	es.BaseAggregate
	_events []es.Event
}

func (this *mockAggregate) HandleMockEvent1(event *mockEvent1) {
	this._events = append(this._events, event)
}

func (this *mockAggregate) HandleMockEvent2(event *mockEvent2) {
	this._events = append(this._events, event)
}

func TestCreateEncoder(t *testing.T) {
	var encoderI Encoder
	factory := es.NewEventFactory()
	encoderI = NewEncoder(factory)
	encoder := encoderI.(*encoder)

	assert.NotNil(t, encoderI, "encoder应该不为空")
	assert.Equal(t, factory, encoder._eventFactory, "factory应该相等")
}

func TestCreateDecoder(t *testing.T) {
	var decoderI Decoder
	factory := es.NewEventFactory()
	decoderI = NewDecoder(factory)
	decoder := decoderI.(*decoder)

	assert.NotNil(t, decoderI, "decoder应该不为空")
	assert.Equal(t, factory, decoder._eventFactory, "factory应该相等")
}

func TestEncoding(t *testing.T) {
	factory := es.NewEventFactory()
	factory.RegisterAggregate(new(mockAggregate))
	encoder, decoder := NewEncoder(factory), NewDecoder(factory)
	in := &mockEvent1{es.WithGuid{es.NewGuid()}, "event1"}

	encodeData, err := encoder.Encode(in, 3)
	assert.NoError(t, err, "事件组包应该没有错误")
	assert.True(t, encodeData != "", "组包内容不为空")

	out, version, err := decoder.Decode(encodeData)
	assert.NoError(t, err, "事件解包应该没有错误")
	assert.Equal(t, 3, version, "跟组包版本不对应")
	assert.Equal(t, out, in, "解包事件跟组包不对应")
}

func TestEncodingNotRegisteredEvent(t *testing.T) {
	factory := es.NewEventFactory()
	factory.RegisterAggregate(new(mockAggregate))
	encoder, decoder := NewEncoder(factory), NewDecoder(factory)
	in := &mockNotRegisterEvent{es.WithGuid{es.NewGuid()}, "event1"}

	encodeData, err := encoder.Encode(in, 3)
	assert.NoError(t, err, "事件组包应该没有错误")
	assert.True(t, encodeData != "", "组包内容不为空")

	_, _, err = decoder.Decode(encodeData)
	assert.Equal(t, es.ErrEventNotRegistered, err, "应该返回事件未注册错误")
}
