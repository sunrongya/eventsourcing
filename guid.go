package eventsourcing

import "github.com/twinj/uuid"

// An item having a GUID
type Guider interface {
	GetGuid() Guid
	SetGuid(Guid)
}

// Base implementation for all Guiders
type WithGuid struct {
	Guid Guid
}

func (e *WithGuid) SetGuid(g Guid) {
	e.Guid = g
}
func (e *WithGuid) GetGuid() Guid {
	return e.Guid
}

type Guid string

// Create a new GUID - use UUID v4
func NewGuid() Guid {
	return Guid(uuid.NewV4().String())
}
