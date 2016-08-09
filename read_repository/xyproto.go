package read_repository

import (
	"github.com/xyproto/pinterface"
)

type XyprotoRepository struct {
	creator pinterface.ICreator
}

func NewXyprotoRepository(creator pinterface.ICreator) *XyprotoRepository {
	return &XyprotoRepository{creator: creator}
}

// Save saves a read model with id to the repository.
func (this *XyprotoRepository) Save(Guid, interface{}) error {
	return nil
}

// Find returns one read model with using an id.
func (this *XyprotoRepository) Find(Guid) (interface{}, error) {
	return nil, nil
}

// FindAll returns all read models in the repository.
func (this *XyprotoRepository) FindAll() ([]interface{}, error) {
	return nil, nil
}

// Remove removes a read model with id from the repository.
func (this *XyprotoRepository) Remove(Guid) error {
	return nil
}
