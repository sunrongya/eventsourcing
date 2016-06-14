package eventsourcing

type Service interface {
	HandleCommands()
	CommandChannel() chan<- Command
	PublishCommand(Command) error
	RestoreAggregate(Guid) Aggregate
}

// common properties for all customer facing services
type service struct {
	commandChannel   chan Command
	store            EventStore
	aggregateFactory func() Aggregate
}

func NewService(store EventStore, aggregateFactory func() Aggregate) Service {
	return &service{
		commandChannel:   make(chan Command),
		store:            store,
		aggregateFactory: aggregateFactory,
	}
}

// Getter for command channel - will allow others to post commands
func (s *service) CommandChannel() chan<- Command {
	return s.commandChannel
}

func (s *service) PublishCommand(command Command) error {
	if err := checkCommand(command); err != nil {
		return err
	}
	s.commandChannel <- command
	return nil
}

// Reads from command channel,
// restores an aggregate,
// processes the command and
// persists received events.
// This method *blocks* until command is available,
// therefore should run in a goroutine
func (s *service) HandleCommands() {
	for {
		c := <-s.commandChannel
		aggregate := s.RestoreAggregate(c.GetGuid())
		s.store.Update(c.GetGuid(), aggregate.Version(), aggregate.ProcessCommand(c))
	}
}

func (s *service) RestoreAggregate(guid Guid) Aggregate {
	if s.aggregateFactory == nil {
		return nil
	}
	aggregate := s.aggregateFactory()
	RestoreAggregate(guid, aggregate, s.store)
	return aggregate
}
