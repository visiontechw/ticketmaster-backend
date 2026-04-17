package domain

import (
	"errors"
	"strings"
)

type EventType struct {
	Base
	Name string `db:"name" gorm:"type:varchar(255);not null"`
}

func (e *EventType) validate() error {

	if strings.TrimSpace(e.Name) == "" {
		return errors.New("event name is required")
	}
	return nil
}

func NewEventType(name string) (*EventType, error) {
	event := &EventType{
		Base: NewBase(),
		Name: strings.TrimSpace(name),
	}

	if err := event.validate(); err != nil {
		return nil, err
	}

	return event, nil
}
