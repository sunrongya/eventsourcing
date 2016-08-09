package eventsourcing

import (
	"fmt"
)

type Service interface {
	HandleCommands()
	CommandChannel() chan<- Command
	PublishCommand(Command) error
	RestoreAggregate(Guid) Aggregate
}

// common properties for all customer facing services
type service struct {
	_commandChannel   chan Command
	_store            EventStore
	_aggregateFactory func() Aggregate
	_commanDispatcher *CommandDispatcher
}

func NewService(store EventStore, aggregateFactory func() Aggregate) Service {
	service := &service{
		_commandChannel:   make(chan Command),
		_store:            store,
		_aggregateFactory: aggregateFactory,
		_commanDispatcher: NewCommandDispatcher(),
	}
	if aggregateFactory != nil {
		service._commanDispatcher.Register(aggregateFactory())
		aggerateEventRegister(aggregateFactory())
	}
	return service
}

// Getter for command channel - will allow others to post commands
func (this *service) CommandChannel() chan<- Command {
	return this._commandChannel
}

func (this *service) PublishCommand(command Command) error {
	if err := checkCommand(command); err != nil {
		return err
	}
	this._commandChannel <- command
	return nil
}

// Reads from command channel,
// restores an aggregate,
// processes the command and
// persists received events.
// This method *blocks* until command is available,
// therefore should run in a goroutine
func (this *service) HandleCommands() {
	for {
		c := <-this._commandChannel
		aggregate := this.RestoreAggregate(c.GetGuid())
		if processCommandFun, ok := this._commanDispatcher.Get(c); ok {
			events := processCommandFun(aggregate)(c)
			for _, event := range events {
				event.SetGuid(c.GetGuid())
			}
			this._store.Update(c.GetGuid(), aggregate.Version(), events)
		} else {
			panic(fmt.Errorf("Unknown command %#v", c))
		}
	}
}

func (this *service) RestoreAggregate(guid Guid) Aggregate {
	if this._aggregateFactory == nil {
		return nil
	}
	aggregate := this._aggregateFactory()
	RestoreAggregate(guid, aggregate, this._store)
	return aggregate
}
