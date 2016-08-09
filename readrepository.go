package eventsourcing

import (
	"errors"
)

// Error returned when a model could not be found.
var ErrCouldNotSaveModel = errors.New("could not save model")

// Error returned when a model could not be found.
var ErrModelNotFound = errors.New("could not find model")

// ReadRepository is a storage for read models.
type ReadRepository interface {
	// Save saves a read model with id to the repository.
	Save(Guid, interface{}) error

	// Find returns one read model with using an id.
	Find(Guid) (interface{}, error)

	// FindAll returns all read models in the repository.
	FindAll() ([]interface{}, error)

	// Remove removes a read model with id from the repository.
	Remove(Guid) error
}

// MemoryReadRepository implements an in memory repository of read models.
type MemoryReadRepository struct {
	_data map[Guid]interface{}
}

// NewMemoryReadRepository creates a new MemoryReadRepository.
func NewMemoryReadRepository() *MemoryReadRepository {
	r := &MemoryReadRepository{
		_data: make(map[Guid]interface{}),
	}
	return r
}

// Save saves a read model with id to the repository.
func (this *MemoryReadRepository) Save(id Guid, model interface{}) error {
	this._data[id] = model
	return nil
}

// Find returns one read model with using an id. Returns
// ErrModelNotFound if no model could be found.
func (this *MemoryReadRepository) Find(id Guid) (interface{}, error) {
	if model, ok := this._data[id]; ok {
		return model, nil
	}

	return nil, ErrModelNotFound
}

// FindAll returns all read models in the repository.
func (this *MemoryReadRepository) FindAll() ([]interface{}, error) {
	models := []interface{}{}
	for _, model := range this._data {
		models = append(models, model)
	}
	return models, nil
}

// Remove removes a read model with id from the repository. Returns
// ErrModelNotFound if no model could be found.
func (this *MemoryReadRepository) Remove(id Guid) error {
	if _, ok := this._data[id]; ok {
		delete(this._data, id)
		return nil
	}

	return ErrModelNotFound
}
