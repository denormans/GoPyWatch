package gopywatch

type EventType uint32

const (
	ProgramDone EventType = 1 << iota
	Restart
	Quit
)

type Event struct {
	Type EventType
}

func NewEvent(eventType EventType) *Event {
	return &Event{
		Type: eventType,
	}
}
