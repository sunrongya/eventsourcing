package estore

import (
	"encoding/json"
	es "github.com/sunrongya/eventsourcing"
)

func NewDecoder(eventFactory *es.EventFactory) Decoder {
	return &decoder{_eventFactory: eventFactory}
}

type decoder struct {
	_eventFactory *es.EventFactory
}

func (this *decoder) Decode(src string) (event es.Event, version int, err error) {
	var r eventRecord
	if err = json.Unmarshal([]byte(src), &r); err != nil {
		return
	}
	version = r.Version
	if event, err = this._eventFactory.GetEvent(r.Type); err != nil {
		return
	}
	err = json.Unmarshal(r.Event, event)

	return
}
