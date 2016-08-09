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

func (this *WithGuid) SetGuid(g Guid) {
	this.Guid = g
}
func (this *WithGuid) GetGuid() Guid {
	return this.Guid
}

type Guid string

// Create a new GUID - use UUID v4
func NewGuid() Guid {
	return Guid(uuid.NewV4().String())
}
