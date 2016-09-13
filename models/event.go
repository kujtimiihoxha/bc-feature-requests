package models

type EventType int

const (
	EVENT_JOIN = iota
	EVENT_LEAVE
	EVENT_MESSAGE
)

type Event struct {
	Type      EventType // JOIN, LEAVE, MESSAGE
	User      string
	Timestamp int // Unix timestamp (secs)
	Content   interface{}
}



