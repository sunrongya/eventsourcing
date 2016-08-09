package eventsourcing


type EventHandler interface {
    HandleEvent(Event)
}


