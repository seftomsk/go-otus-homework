package storage

import "time"

type Event struct {
	id           string
	userID       string
	title        string
	description  EventDescription
	datetime     time.Time
	duration     time.Duration
	remindBefore time.Duration
}

type EventDescription struct {
	Data  string
	Valid bool
}

func (e *Event) ID() string {
	return e.id
}

func (e *Event) UserID() string {
	return e.userID
}

func (e *Event) Title() string {
	return e.title
}

func (e *Event) Datetime() time.Time {
	return e.datetime
}

func (e *Event) Duration() time.Duration {
	return e.duration
}

func (e *Event) Description() EventDescription {
	return e.description
}

func (e *Event) RemindBefore() time.Duration {
	return e.remindBefore
}

func NewEvent(
	id string,
	userID string,
	title string,
	description EventDescription,
	datetime time.Time,
	duration time.Duration,
	remindBefore time.Duration) *Event {
	return &Event{
		id:           id,
		userID:       userID,
		title:        title,
		description:  description,
		datetime:     datetime.UTC(),
		duration:     duration,
		remindBefore: remindBefore,
	}
}
