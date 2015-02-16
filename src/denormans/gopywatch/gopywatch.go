package gopywatch

type EventType uint32

const (
	ProgramStarted EventType = 1 << iota
	ProgramDone
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
