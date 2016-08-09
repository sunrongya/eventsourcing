package estore

import (
	"encoding/json"
	"time"

	es "github.com/sunrongya/eventsourcing"
)

func NewEncoder(eventFactory *es.EventFactory) Encoder {
	return &encoder{_eventFactory: eventFactory}
}

type encoder struct {
	_eventFactory *es.EventFactory
}

func (this *encoder) Encode(event es.Event, version int) (string, error) {
	// TODO 添加事件是否已注册判断
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	r := &eventRecord{
		Type:      this._eventFactory.EventStringType(event),
		Version:   version,
		Timestamp: time.Now(),
		Event:     json.RawMessage(eventBytes),
	}

	b, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
