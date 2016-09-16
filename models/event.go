package models

type EventType int

const (
	EVENT_JOIN = iota
	EVENT_LEAVE
	EVENT_MESSAGE
)

type Event struct {
	Type      EventType   `json:"type"` // JOIN, LEAVE, MESSAGE
	User      string      `json:"user"`
	Timestamp int         `json:"timestamp"` // Unix timestamp (secs)
	Content   interface{} `json:"content"`
}
