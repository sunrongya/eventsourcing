package estore

import (
	"encoding/json"
	"fmt"
	es "github.com/sunrongya/eventsourcing"
	"github.com/xyproto/pinterface"
	"math"
	"time"
)

var ErrEventEncode = fmt.Errorf("event encode error")
var ErrEventDecode = fmt.Errorf("event decode error")
var ErrNoEventsFound = fmt.Errorf("could not find events")
var ErrLoadEvents = fmt.Errorf("load events error")

type Encoder interface {
	Encode(es.Event, int) (string, error)
}

type Decoder interface {
	Decode(string) (es.Event, int, error)
}

type XyprotoEStore struct {
	es.EventChannel
	_creator pinterface.ICreator
	_encoder Encoder
	_decoder Decoder
}

func NewXyprotoEStore(creator pinterface.ICreator, encoder Encoder, decoder Decoder) *XyprotoEStore {
	return &XyprotoEStore{
		EventChannel: es.NewEventChannel(),
		_creator:     creator,
		_encoder:     encoder,
		_decoder:     decoder,
	}
}

func (this *XyprotoEStore) Find(guid es.Guid) (events []es.Event, version int) {
	list, err := this._creator.NewList(this.listName(guid))
	if err != nil {
		return nil, -1
	}
	allEvents, err := list.GetAll()
	if err != nil || len(allEvents) == 0 {
		return nil, 0
	}

	return this.getEvents(allEvents)
}

func (this *XyprotoEStore) Update(guid es.Guid, version int, events []es.Event) error {
	_, lastVersion, err := this.getLastEvent(guid)
	if err != nil {
		return err
	}
	if lastVersion != version {
		return fmt.Errorf("Optimistic locking exeption - client has version %v, but store %v", version, lastVersion)
	}

	list, err := this._creator.NewList(this.listName(guid))
	if err != nil {
		return err
	}

	for i, event := range events {
		event.SetGuid(guid)
		eventstr, err := this.encode(event, lastVersion+i+1)
		if err != nil {
			return err
		}
		if err := list.Add(eventstr); err != nil {
			return nil
		}
	}
	this.AppendEvents(events)

	return nil
}

// 用这个方法时要注意：返回的内容是最后N条事件(N := Min(offset+batchSize, version))
func (this *XyprotoEStore) GetEvents(guid es.Guid, offset int, batchSize int) []es.Event {
	_, version, err := this.getLastEvent(guid)
	if err != nil {
		return nil
	}
	list, err := this._creator.NewList(this.listName(guid))
	if err != nil {
		return nil
	}
	until := int(math.Min(float64(offset+batchSize), float64(version)))

	lastEvents, err := list.GetLastN(until)
	if err != nil || len(lastEvents) == 0 {
		return nil
	}
	result, _ := this.getEvents(lastEvents)

	return result
}

func (this *XyprotoEStore) getLastEvent(guid es.Guid) (event es.Event, version int, err error) {
	list, err := this._creator.NewList(this.listName(guid))
	if err != nil {
		return event, -1, err
	}
	eventStr, err := list.GetLast()
	if err != nil {
		return event, -1, err
	}
	if eventStr == "" {
		return event, 0, nil
	}
	return this.decode(eventStr)
}

func (this *XyprotoEStore) getEvents(rawEvents []string) (events []es.Event, version int) {
	for _, strEvent := range rawEvents {
		event, tVersion, err := this.decode(strEvent)
		if err != nil {
			panic(err.Error())
		}
		events = append(events, event)
		version = tVersion
	}
	return
}

func (this *XyprotoEStore) encode(event es.Event, version int) (string, error) {
	return this._encoder.Encode(event, version)
}

func (this *XyprotoEStore) decode(eventStr string) (event es.Event, version int, err error) {
	return this._decoder.Decode(eventStr)
}

func (this *XyprotoEStore) listName(guid es.Guid) string {
	return "Event" + string(guid)
}

type eventRecord struct {
	Type      string
	Version   int
	Timestamp time.Time
	Event     json.RawMessage
}
